package textbuilder

import (
	"fmt"
	"strings"
)

func BuildPrefix(prefix TextGetter) Transformer {
	return func(content TextGetter) TextGetter {
		return &prefixImpl{prefix, content}
	}
}

func BuildSuffix(suffix TextGetter) Transformer {
	return func(content TextGetter) TextGetter {
		return &suffixImpl{content, suffix}
	}
}

func BuildSurrounding(prefix TextGetter, suffix TextGetter) Transformer {
	return func(content TextGetter) TextGetter {
		return &surroundingImpl{prefix, content, suffix}
	}
}

func BuildSameSurrounding(s TextGetter) Transformer {
	return BuildSurrounding(s, s)
}

func BuildRepeat(times int) Breeder {
	return func(content TextGetter) TextList {
		return Repeat(content, times)
	}
}

func BuildRepeatByLen(v interface{}) Breeder {
	return func(content TextGetter) TextList {
		return RepeatByLen(content, v)
	}
}

func BuildJoin(separator TextGetter) Distiller {
	return func(list TextList) TextGetter {
		return &joinImpl{separator, list}
	}
}

type prefixImpl struct {
	prefixImpl TextGetter
	content    TextGetter
}

func (p *prefixImpl) String() string {
	content := p.content.String()
	if content != "" {
		return p.prefixImpl.String() + content
	}

	return ""
}
func (p *prefixImpl) Post() PostProcessor {
	return NewPost(p)
}

type suffixImpl struct {
	content TextGetter
	suffix  TextGetter
}

func (s *suffixImpl) String() string {
	content := s.content.String()
	if content != "" {
		return content + s.suffix.String()
	}

	return ""
}
func (s *suffixImpl) Post() PostProcessor {
	return NewPost(s)
}

type surroundingImpl struct {
	prefix  TextGetter
	content TextGetter
	suffix  TextGetter
}

func (s *surroundingImpl) String() string {
	content := s.content.String()
	if content != "" {
		return fmt.Sprintf(
			"%s%s%s",
			s.prefix.String(), s.content.String(), s.suffix.String(),
		)
	}

	return ""
}
func (s *surroundingImpl) Post() PostProcessor {
	return NewPost(s)
}

type joinImpl struct {
	separator TextGetter
	getters   TextList
}

func (j *joinImpl) String() string {
	resultOfGetters := make([]string, 0)

	for i := 0; i < j.getters.Len(); i++ {
		getterResult := j.getters.Get(i).String()
		if getterResult == "" {
			continue
		}

		resultOfGetters = append(resultOfGetters, getterResult)
	}

	return strings.Join(resultOfGetters, j.separator.String())
}
func (j *joinImpl) Post() PostProcessor {
	return NewPost(j)
}
