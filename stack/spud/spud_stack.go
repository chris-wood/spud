package spud

import "log"
import "io/ioutil"

import "crypto/x509"
import "encoding/pem"
import "encoding/json"

import tlvCodec "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/tables/lpm"
import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/stack"
import "github.com/chris-wood/spud/stack/config"
import "github.com/chris-wood/spud/stack/cache"
import "github.com/chris-wood/spud/stack/pit"
import "github.com/chris-wood/spud/stack/component"
import "github.com/chris-wood/spud/stack/component/transport"
import "github.com/chris-wood/spud/stack/component/tunnel"
import "github.com/chris-wood/spud/stack/component/connector"
import "github.com/chris-wood/spud/stack/component/codec"
import "github.com/chris-wood/spud/stack/component/crypto"
import "github.com/chris-wood/spud/stack/component/crypto/validator"
import "github.com/chris-wood/spud/stack/component/crypto/context"

type SpudStack struct {
	cryptoComponent component.Component
	codecComponent  component.Component
    transportComponent  component.Component
    tunnelComponent *tunnel.TunnelComponent
	head            component.Component
	bottom          component.Component

	pendingMap   map[string]stack.MessageCallback
	pendingQueue chan *messages.MessageWrapper
	prefixTable  lpm.LPM
}

func (s *SpudStack) Enqueue(msg *messages.MessageWrapper) {
	s.head.Enqueue(msg)
}

func (s *SpudStack) Dequeue() *messages.MessageWrapper {
	return <-s.pendingQueue
}

func (s *SpudStack) Cancel(msg *messages.MessageWrapper) {
	_, ok := s.pendingMap[msg.Identifier()]
	if ok {
		delete(s.pendingMap, msg.Identifier())
	}
}

func (s *SpudStack) Get(msg *messages.MessageWrapper, callback stack.MessageCallback) {
	if msg.GetPacketType() != tlvCodec.T_INTEREST {
		return
	}

	// Register the callback
	s.pendingMap[msg.Identifier()] = callback

	// Enqueue the message into the top of the stack
	s.head.Enqueue(msg)
}

func (s *SpudStack) Service(prefix *name.Name, callback stack.MessageCallback) {
    if prefix == nil {
        return // XXX: return an error
    }
	nameComponents := prefix.SegmentStrings()
	s.prefixTable.Insert(nameComponents, callback)
}

func (s *SpudStack) processInputQueue() {
	for {
		// Dequeue messages as they arrive
		msg := s.head.Dequeue()

		switch msg.GetPacketType() {
		case tlvCodec.T_INTEREST:
			requestName := msg.Name()
			numSsegments := len(requestName.Segments)

			handled := false
			for index := 1; index <= numSsegments; index++ {
				nameComponents := requestName.SegmentStrings()
				callbackInterface, ok := s.prefixTable.Lookup(nameComponents)

				if ok {
					callback := callbackInterface.(stack.MessageCallback)
					callback(msg)
					s.prefixTable.Drop(nameComponents)
					handled = true
				}
			}

			if !handled {
				s.pendingQueue <- msg
			}
		default:
			callback, ok := s.pendingMap[msg.Identifier()]
			if ok {
				callback(msg)
			} else {
				s.pendingQueue <- msg
			}
		}
	}
}

func (s *SpudStack) AddSession(session *tunnel.Session, baseName *name.Name) {
    s.tunnelComponent.AddSession(session, baseName)
}

func Create(config config.StackConfig) (*SpudStack, error) {
	var err error

	// Create the link
	var fc connector.ForwarderConnector
	switch config.Link {
	case "tcp":
		locator := config.FwdAddress
		fc, err = connector.NewAthenaTCPForwarderConnector(locator)
		if err != nil {
			panic(err)
		}
		break
	case "loopback":
		fallthrough
	default:
		fc, _ = connector.NewLoopbackForwarderConnector()
		break
	}

	if fc == nil {
		log.Panic("Catastrophic failure: a connector was not created.")
	}

	// Create the shared data structures
	stackCache := cache.NewCache()
	stackPit := pit.NewPIT()

	// Create codec
	codecComponent := codec.NewCodecComponent(fc, stackCache, stackPit)

    // Create the tunnel component
    tunnelComponent := tunnel.NewTunnelComponent(codecComponent)

    // Create the transport component
    transportComponent := transport.NewTransportComponent(tunnelComponent)

	// Create crypto component
	// XXX: delegate validator and encryptor creation to the crypto component
	trustStore := context.NewTrustStore()
	cryptoComponent := crypto.NewCryptoComponent(trustStore, transportComponent)

	// Create the right crypto processors
	if len(config.Keys) > 0 {
		for _, keyFileName := range config.Keys {
			keyData, err := ioutil.ReadFile(keyFileName)
			block, _ := pem.Decode(keyData)
			privateKey, parseError := x509.ParsePKCS1PrivateKey(block.Bytes)
			if parseError != nil {
				log.Printf("Failed to parse private key: %s", err)
			} else {
				rsaProcessor, _ := validator.NewRSAProcessorWithKey(privateKey)
				cryptoComponent.AddCryptoProcessor("/", rsaProcessor)
			}
		}
	} else {
		rsaProcessor, _ := validator.NewRSAProcessor(2048)
		cryptoComponent.AddCryptoProcessor("/", rsaProcessor)
	}

	// Assemble the stack
	stack := &SpudStack{
		cryptoComponent: cryptoComponent,
		codecComponent:  codecComponent,
        transportComponent: transportComponent,
        tunnelComponent: tunnelComponent,
		head:            cryptoComponent,
		pendingQueue:    make(chan *messages.MessageWrapper, config.PendingBufferSize), // random constant -- make this configurable
		pendingMap:      make(map[string]stack.MessageCallback),
		prefixTable:     &lpm.StandardLPM{},
	}

	// Start the stack processing loop -- the order here matters
	go codecComponent.ProcessEgressMessages()
	go codecComponent.ProcessIngressMessages()
    	go tunnelComponent.ProcessEgressMessages()
	go tunnelComponent.ProcessIngressMessages()
    	go transportComponent.ProcessEgressMessages()
	go transportComponent.ProcessIngressMessages()
	go cryptoComponent.ProcessEgressMessages()
	go cryptoComponent.ProcessIngressMessages()
	go stack.processInputQueue()

	return stack, nil
}

func CreateRaw(configString string) (*SpudStack, error) {
	var configStruct config.StackConfig
	err := json.Unmarshal([]byte(configString), &configStruct)
	if err == nil {
		return Create(configStruct)
	}
	return nil, err
}

func CreateTest() (*SpudStack, error) {
	return CreateRaw(`{"link": "loopback"}`)
}
