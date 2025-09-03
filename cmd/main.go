package main

func main() {
	ctx := ctrl.SetupSignalHandler()

	operator, err := NewMyOperator()
	if err != nil {
		panic(err)
	}

	if err := operator.Start(ctx); err != nil {
		panic(err)
	}
}
