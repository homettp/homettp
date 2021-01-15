package forms

type Bag map[string][]string

func (b Bag) Add(field, message string) {
	b[field] = append(b[field], message)
}
