package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/badoux/checkmail"
	"github.com/paemuri/brdoc"
	"golang.org/x/crypto/bcrypt"
)

type validator interface {
	CheckIsEmpty(input interface{}, nameInput string) error
	CheckLen(size int, input interface{}) error
	CheckEmail(email string) error
	CheckPassword(password string) (string, error)
	CheckHash(pswWithHash, password string) error
	CheckCPFOrCNPJ(cpfOrCnpj string) error
	FormatedInput(input string) (string, error)
}

type validadorObject struct{}

func (v *validadorObject) CheckIsEmpty(input interface{}, nameInput string) error {
	if reflect.ValueOf(input).Kind() == reflect.String {
		str, ok := input.(string)
		if !ok || str == "" {
			return errors.New(fmt.Sprintf("%v é obrigatório e não pode está em branco", nameInput))
		}
	}
	if reflect.ValueOf(input).Kind() == reflect.Int {
		_, ok := input.(int)
		if !ok {
			return errors.New(fmt.Sprintf("%v é obrigatório e não pode está em branco", nameInput))
		}
	}
	if reflect.ValueOf(input).Kind() == reflect.Bool {
		_, ok := input.(bool)
		if !ok {
			return errors.New(fmt.Sprintf("%v é obrigatório e não pode está em branco", nameInput))
		}
	}

	return nil
}
func (v *validadorObject) CheckLen(size int, input interface{}) error {
	if size < len(input.(string)) {
		return errors.New("O tamanho excede o permitido")
	}
	return nil
}
func (v *validadorObject) CheckEmail(email string) error {
	if err := checkmail.ValidateFormat(email); err != nil {
		return errors.New("O e-mail inserido é invalido")
	}
	return nil
}
func (v *validadorObject) CheckPassword(password string) (string, error) {
	hashByte, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashByte), err
}
func (v *validadorObject) CheckHash(pswWithHash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(pswWithHash), []byte(password))
}
func (v *validadorObject) CheckCPFOrCNPJ(cpfOrCnpj string) error {
	okCpf := brdoc.IsCPF(cpfOrCnpj)
	okCnpj := brdoc.IsCNPJ(cpfOrCnpj)
	if okCpf || okCnpj {
		return nil
	}
	return errors.New("CPF ou CNPJ invalidos")
}
func (v *validadorObject) FormatedInput(input string) (string, error) {
	return strings.TrimSpace(input), nil
}

func NewValidator() validator {
	return &validadorObject{}
}
