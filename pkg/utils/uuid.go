package utils

import (
	"github.com/google/uuid"
	"strconv"
	"time"
)

/*
*
uuid.EnableRandPool

	启用随机数池 (EnableRandPool = true)：当启用随机数池时，UUID库会在初始化时填充一个随机数池，然后在生成UUID时，
		会优先使用随机数池中的字节。这样做的好处是，避免了频繁地从操作系统获取随机数，从而提高了性能。
		由于随机数池中的随机字节是提前生成的，生成UUID时可以直接从池中获取，减少了随机数生成的开销。
		启用随机数池适用于需要高性能的场景，特别是在大量生成UUID的情况下。

	不启用随机数池 (EnableRandPool = false)：如果不启用随机数池，那么每次生成UUID时都会直接调用操作系统的随机数生成器，
		获取新的随机字节。这样做的好处是保证了更高的随机性，但可能会对性能产生一些影响，因为频繁调用操作系统的随机数生成器可能会产生一些开销。
		不启用随机数池适用于对随机性要求更高、而不是特别关注性能的场景。
	总之，EnableRandPool 这个配置选项允许您在Google为Go语言实现的UUID库中权衡性能和随机性。
		如果您的应用需要在高性能的情况下生成大量UUID，可以考虑启用随机数池。
		如果您更关注随机性，而不太关心性能，可以选择不启用随机数池
*/
func init() {
	uuid.EnableRandPool()
}

// UUID4 使用Google的uuid库对版本4的实现
//
//	如果有错误就返回当前时间的 ns时间戳字符串(但是基本不会有错误的)
func UUID4() string {
	id, err := uuid.NewRandom()
	if err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	return id.String()
}
