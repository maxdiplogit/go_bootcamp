package structs

type Example struct {
    A string `json:"a"`             // rename to "a"
    B string `json:"b,omitempty"`   // rename and omit if empty
    C string `json:",omitempty"`    // keep name "C", but omit if empty
    D string `json:"-"`             // never serialize
    E string `json:"-,"`            // literally serialize as "-" (weird edge case)
    F string                         // no tag — serialize as "F"
}
