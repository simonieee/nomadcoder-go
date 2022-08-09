package mydict

import "errors"

// Dictionary type
type Dictionary map[string]string

var (
	errNotFound = errors.New("Not Found")
 	errCantUpdate = errors.New("Cant update non-existing word")
 	errWordExists = errors.New("That word already exists")
)

func (d Dictionary) Search(word string) (string, error){
	value, exists := d[word]
	if exists {
		return value, nil
	}
	return "", errNotFound
}

func (d Dictionary) Add(word, def string) error {
	_, err := d.Search(word)
	switch err {
	case errNotFound:
		d[word] = def
	case nil:
		return errWordExists 
	}
	return nil
}

// Update a word
func (d Dictionary) Update(word, definition string) error{
	_, err := d.Search(word)
	switch err {
	case nil:
		d[word] = definition
	case errNotFound:
		return errCantUpdate
	}
	return nil
}

// Delete a word
func (d Dictionary) Delete(word string) {
	delete(d, word)
}