package data

func getMessages() map[string]string {
	return map[string]string{
		"startMessage":         "👋Приветствуем!\n Бот поможет провести Вам розыгрыш на канале.",
		"addTargetChannels":    "Введите id канала, в котором будет опубликован розыгрыш. Вы должны быть администратором этого канала и добавить в администраторы бота.\\n Чтобы узнать id приватного канала, перешлите любой пост из этого канала специальному боту @getmyid_bot, нижняя строка-ответ от бота и будет являться id Вашего канала. Перешлите этот id боту",
		"successRefreshAdmins": "Список админов успешно обновлен",
	}
}
