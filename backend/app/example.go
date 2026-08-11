package app

// var pp = map[string]string{
// 	"1": "2",
// 	"c": "b",
// }

//申明并绑定对列和交换机
// 	rabbitmq.Declare("e.haha", "q.haha")

//发送Rabbitmq消息
// 	for i := 0; i < 50; i++ {
// 		rabbitmq.Send("e.haha", "q.haha", pp)
// 		fmt.Println(i)
// 	}

//接受Rabbitmq消息
// rabbitmq.Receive("q.haha", func(s string) error {
// 	fmt.Println(s)
// 	return nil
// })

//发送kafka消息
// for i := 0; i < 100; i++ {
// 	kafka.Send("tt1", pp)
// }

// 接受kafka消息
// kafka.Receive("tt1", func(s string) error {
// 	fmt.Println(s)
// 	return nil
// })

//mongodb插入单条数据
// var map1 = map[string]any{
// 	"a": "a1",
// 	"b": "b1",
// }
// datasource.GetMgoDb().Collection("nt_push_record").InsertOne(context.Background(), map1)

//mongodb根据条件查询数据
// filter := bson.D{{}}
// cursor, err := datasource.GetMgoDb().Collection("nt_push_record").Find(nil, filter)
// for cursor.Next(context.Background()) {
// 	var aa any
// 	cursor.Decode(&aa)
// 	fmt.Println(aa)
// }

//定时任务
// schedule.AddJob("订单超时", "*/5 * * * * ?", func() {
// 	fmt.Println("now:" + time.Now().Format("2006-01-02 15:04:05"))
// })
