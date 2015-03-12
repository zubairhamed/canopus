package goap

/*
	// Discovery Test
	wg.Add(1)
	client = newTestClient()
	client.OnSuccess(func(msg *Message) {
		defer wg.Done()

		s := PayloadAsString(msg.Payload)

		if s != "</testGET>,</discoveryService1>,</discoveryService2>,</discoveryService3>," {
			t.Error("Unexpected return for discovery service")
		}
	})

	msg = NewMessageOfType(TYPE_CONFIRMABLE, 23456)
	msg.Code = GET
	msg.AddOptions(NewPathOptions("/.well-known/core"))
	msg.Token = []byte(GenerateToken(8))

	go client.Send(msg)

	wg.Wait()
	client.Close()
*/
