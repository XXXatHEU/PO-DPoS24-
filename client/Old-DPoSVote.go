package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

//投票策略
func Old_votStrategy(data dposJsonStruct) map[int][]int {
	//全部的信誉值
	SelectVersion = data.IntData
	Old_PersonVoteList := data.Old_PersonVoteList
	Old_GrupMclic := data.Old_MclicNodeSGrup //所有的恶意节点
	//myReputatil := data.ReputationDetail
	//var resultArry map[int][]int
	//全局恶意节点map
	//下标对应
	fmt.Println("MyInitNum", MyInitNum, "  data.GroupPersons", data.GroupPersons)

	GrupMclic = data.MclicNodeSGrup
	for index, num := range GrupMclic[data.GroupNum] {
		inMclieGrup[num] = index
	}
	for _, num := range data.GroupPersons {
		inGrup[num] = true
	}
	/*
		如果是恶意节点
		 	作恶就是对诚实节点投反对票			对恶意节点投非常赞同
			诚实就是对诚实节点投赞同或者弃权	对恶意节点投反对
		如果是诚实节点
			随机投票
		（要有一些概率）
	*/
	//fmt.Println("MyInitNum", MyInitNum, "   ", Old_PersonVoteList)
	recognition_probability += reconition
	for nodeNum, _ := range Old_PersonVoteList {
		Old_PersonVoteList[nodeNum][0] = 0
		Old_PersonVoteList[nodeNum][1] = 0
		Probability := data.ReputationDetail.Extra.MmaliciousnessProbability
		temptype := nodetype
		// 设置随机数种子
		rand.Seed(time.Now().UnixNano())
		// 生成随机数
		randomNumber := rand.Float64()
		//如果我是恶意节点 在作恶概率下 有可能这次不选择作恶，那么在这决定是否作恶
		if nodetype == unhealthyNode && randomNumber > Probability { //不作恶
			temptype = healthyNode
		}

		if temptype == unhealthyNode { //如果我是恶意节点
			//对于诚实节点 投反对票
			if _, ok := Old_GrupMclic[nodeNum]; ok { //这是恶意节点
				Old_PersonVoteList[nodeNum][0] = 1
			} else { //对于诚实节点投赞成票
				Old_PersonVoteList[nodeNum][1] = 1
			}
		} else { //如果我是诚实节点

			rand.Seed(time.Now().UnixNano())
			randomNumber := rand.Float64()
			//诚实节点会逐渐提升自己对恶意节点的认知
			if _, ok := Old_GrupMclic[nodeNum]; ok && randomNumber <= recognition_probability {
				Old_PersonVoteList[nodeNum][1] = 1 //对于恶意节点投非常反对
			} else {
				rand.Seed(time.Now().UnixNano())
				// 生成0到3范围内的随机数
				randomNumber := rand.Intn(2)
				Old_PersonVoteList[nodeNum][randomNumber] = 1
			}
		}
	}
	fmt.Println("Old_PersonVoteList", Old_PersonVoteList)
	return Old_PersonVoteList
}

func Old_SendVoteResult(Old_PersonVoteList map[int][]int, getdata dposJsonStruct) {
	Log.Info("让我进行投票,MyInitNum:", MyInitNum)
	var outmsg dposJsonStruct

	var GetValue = make(map[int]float64) //恶意交易

	ValidtxMapMutex.Lock()
	for key, value := range HealthTxMap {
		if HealthTxNum == 0 {
			break
		}
		GetValue[key] += float64(value) / float64(HealthTxNum) * 20
	}
	for key, value := range MalicTxMap {
		if HealthTxNum == 0 {
			break
		}
		GetValue[key] -= float64(value) / float64(HealthTxNum) * 35
	}
	ValidtxMapMutex.Unlock()

	outmsg.Comm = Old_ReplyVote
	outmsg.Old_GroupReplyVote = Old_PersonVoteList
	outmsg.InitNodeNum = MyInitNum
	//还未使用
	outmsg.TokenChanges = GetValue
	// 将结构体转换为 JSON 字节流
	jsonBytes, err1 := json.Marshal(outmsg)
	if err1 != nil {
		fmt.Println("转换失败:", err1)
		return
	}
	connWrite(DNSConn, []byte(jsonBytes))
	Log.Info("让我进行投票完毕,MyInitNum:", MyInitNum)
}
