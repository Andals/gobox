/**
* @file base.go
* @brief format msg before send to writer
* @author ligang
* @date 2016-07-12
 */

package formater

type IFormater interface {
	Format(level int, msg []byte) []byte
}
