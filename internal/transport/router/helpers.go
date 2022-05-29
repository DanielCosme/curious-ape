package router

// // The background() helper accepts an arbitrary function as a parameter.
// func (a *main.application) background(fn func()) {
// 	// Launch a background goroutine.
// 	go func() {
// 		// Recover any panic.
// 		defer func() {
// 			if err := recover(); err != nil {
// 				//a.logger.PrintError(fmt.Errorf("%s", err), nil)
// 				a.logger.Println(err)
// 			}
// 		}()
//
// 		// Execute the arbitrary function that we passed as the parameter.
// 		fn()
// 	}()
// }
//
