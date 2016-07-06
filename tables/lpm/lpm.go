package lpm

type prefixKey {
    value string
}

func (pk prefixKey) Compare() bool {

}

type prefixTable struct {
    fixedKeys []prefixKey
    regexKeys []prefixKey
}

type LPM struct {
    tables []prefixTable
}

func (l LPM) Insert(key string, value interface{}) {

}

func (l LPM) Lookup(key string, value interface{}) {

}
