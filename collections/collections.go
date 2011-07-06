package collections

type (
	Any interface{}
	DoHandler func(Any)bool
	Collection interface {
		Do(func(Any)bool)
	}
)

func GetRange(c Collection, start, length int) []Any {
	end := start + length
	items := make([]Any, length)
	i := 0
	j := 0
	c.Do(func(item Any)bool{
		if i >= start {
			if i < end {
				items[j] = item
				j++
			} else {
				return false
			}
		}
		i++
		return true
	})
	return items[:j]
}
