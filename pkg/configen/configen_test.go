package configen

import (
	"testing"
)

func TestInsertYamlbyTxtstatement(t *testing.T) {
	InsertYamlbyTxtstatement("..\\..\\files\\SCFile", "..\\..\\files\\source.yaml", "..\\..\\files\\out.yaml")

}
