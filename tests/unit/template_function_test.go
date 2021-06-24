package unit_test

import (
	"bytes"
	"github.com/Masterminds/sprig/v3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strconv"
	txtTemplate "text/template"
)

var _ = Describe("Helm template function verification", func() {

	It("checks valid JVM formats", func() {

		var v, e = checkTemplate(buildJvmRegexTemplate("1000m"))
		Expect(e).To(BeNil())
		Expect(v).To(BeTrue())

		v, e = checkTemplate(buildJvmRegexTemplate("1000M"))
		Expect(e).To(BeNil())
		Expect(v).To(BeTrue())

		v, e = checkTemplate(buildJvmRegexTemplate("110.2k"))
		Expect(e).To(BeNil())
		Expect(v).To(BeTrue())

		v, e = checkTemplate(buildJvmRegexTemplate("110.2K"))
		Expect(e).To(BeNil())
		Expect(v).To(BeTrue())

		v, e = checkTemplate(buildJvmRegexTemplate("0.1g"))
		Expect(e).To(BeNil())
		Expect(v).To(BeTrue())

		v, e = checkTemplate(buildJvmRegexTemplate("0.1G"))
		Expect(e).To(BeNil())
		Expect(v).To(BeTrue())

		v, e = checkTemplate(buildJvmRegexTemplate("0.1e"))
		Expect(e).To(BeNil())
		Expect(v).To(BeTrue())

		v, e = checkTemplate(buildJvmRegexTemplate("0.1E"))
		Expect(e).To(BeNil())
		Expect(v).To(BeTrue())

		v, e = checkTemplate(buildJvmRegexTemplate("0.1p"))
		Expect(e).To(BeNil())
		Expect(v).To(BeTrue())

		v, e = checkTemplate(buildJvmRegexTemplate("0.1P"))
		Expect(e).To(BeNil())
		Expect(v).To(BeTrue())

		v, e = checkTemplate(buildJvmRegexTemplate("0.1t"))
		Expect(e).To(BeNil())
		Expect(v).To(BeTrue())

		v, e = checkTemplate(buildJvmRegexTemplate("0.1T"))
		Expect(e).To(BeNil())
		Expect(v).To(BeTrue())
	})

	It("checks invalid JVM formats", func() {

		var v, e = checkTemplate(buildJvmRegexTemplate(".2k"))
		Expect(e).To(BeNil())
		Expect(v).To(BeFalse())

		v, e = checkTemplate(buildJvmRegexTemplate("1024foomanchoo"))
		Expect(e).To(BeNil())
		Expect(v).To(BeFalse())

		v, e = checkTemplate(buildJvmRegexTemplate("1024 M"))
		Expect(e).To(BeNil())
		Expect(v).To(BeFalse())

		v, e = checkTemplate(buildJvmRegexTemplate("abcdM"))
		Expect(e).To(BeNil())
		Expect(v).To(BeFalse())

		v, e = checkTemplate(buildJvmRegexTemplate("1024 G"))
		Expect(e).To(BeNil())
		Expect(v).To(BeFalse())

		v, e = checkTemplate(buildJvmRegexTemplate("1024"))
		Expect(e).To(BeNil())
		Expect(v).To(BeFalse())
	})
})

func buildJvmRegexTemplate(value string) (string, map[string]string) {
	Expect(value).ToNot(BeEmpty())
	jvmFormatRegex := `^(([0]\.\d*)+|(^[1-9]\d*)+(\.\d+)?)(?i)(k|m|g|e|p|t){1}$`
	return `{{ regexMatch .exp .input }}`, map[string]string{"exp": jvmFormatRegex, "input": value}
}

func checkTemplate(tpl string, vars interface{}) (bool, error) {
	t := txtTemplate.Must(txtTemplate.New("test-tpl-func-").Funcs(sprig.TxtFuncMap()).Parse(tpl))
	var b bytes.Buffer
	err := t.Execute(&b, vars)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(b.String())
}
