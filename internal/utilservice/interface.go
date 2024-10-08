// utilservice package implements functions that process the matched order list.
package utilservice

type UtilInterface interface {
	// init initialize the pair order list.
	init()

	// start executes the pair trading server.
	start()

	// close signals the service to terminate
	close()
}
