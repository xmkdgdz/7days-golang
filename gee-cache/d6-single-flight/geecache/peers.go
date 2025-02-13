/* 
    流程2：远程节点获取数据（分布式）

    使用一致性哈希选择节点        是                                    是
        |-----> 是否是远程节点 -----> HTTP 客户端访问远程节点 --> 成功？-----> 服务端返回返回值
                        |  否                                    ↓  否
                        |----------------------------> 回退到本地节点处理。
 */

package geecache

// PeerPicker interface 根据传入的 key 选择相应节点 PeerGetter
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter interface 从对应 group 查找缓存值，类似 HTTP 客户端
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}