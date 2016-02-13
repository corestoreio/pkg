package phpsession

import "github.com/corestoreio/csfw/util/php/phpserialize"

const SEPARATOR_VALUE_NAME rune = '|'

type PhpSession map[string]phpserialize.PhpValue
