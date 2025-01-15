package deliveryController

func AddDelivery(requestBody RequestCreateDelivery) int {

	if err := insertDelivery(requestBody); err != nil {
		return 400
	}

	return 200
}
