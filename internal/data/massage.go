package data

const (
	START_message                       = "👋Приветствуем!\nБот поможет Вам провести розыгрыш на канале."
	NO_ADMINS_message                   = "Нет админов для обновления"
	NO_GIVES_message                    = "У тебя нет конкурсов"
	SUCCESS_REFRESH_ADMINS_message      = "Список админов успешно обновлен"
	ENTER_GIVE_TITLE_message            = "Введи название конкурса"
	ENTER_GIVE_DESCRIPTION_message      = "Введите описание конкурса"
	UPLOAD_GIVE_IMAGE_message           = "Загрузи картинку конкурса. Не более 50 Мб"
	ENTER_GIVE_START_FINISH_message     = "Введи период проведения конкурса.\nСтарт - Финиш\n(ДД.ММ.ГГГГ ЧЧ:ММ - ДД.ММ.ГГГГ ЧЧ:ММ)\nвремя указывать по МСК\nНапример:\n01.09.2022 12:00 - 31.09.2022 20:00"
	ENTER_WINNERS_COUNT_message         = "Введи колличество победителей"
	SELECT_OWN_GIVE_message             = "Вот все твои конкурсы. Выбирай"
	ENTER_TARGET_CHANNEL_message        = "Введи id канала, в котором будет опубликован розыгрыш. Вы должны быть администратором этого канала и добавить в администраторы бота.\nЧтобы узнать id канала, перешлите любой пост из этого канала специальному боту @getmyid_bot, нижняя строка-ответ от бота и будет являться id Вашего канала. Перешлите этот id боту"
	ENTER_SUBSCRIPTION_CHANNELS_message = "Введи id каналов на котовые нужно проверить подписку (один или несколько через пробел). Вы должны добавить в администраторы бота.\nЧтобы узнать id канала, перешлите любой пост из этого канала специальному боту @getmyid_bot, нижняя строка-ответ от бота и будет являться id Вашего канала."
	GIVE_OUTPUT_message                 = "_Вот твой конкурс_\n\n*Канал:* %s\n*Проверка подписки:* %s\n*Колличество победителей*: %d\n*Старт:* %s\n*Финиш:* %s\n*Статус:* %s"
	GIVE_CONTENT_message                = "*%s*\n\n%s"
	GIVE_SUCCESSFULLY_ACTIVATE_message  = "Конкурс успешно активирован и будет опубликован в соответствии с параметрами"
	GIVE_SUCCESSFULL_DEACTIVATE_message = "Конкурс успешно снят с публикации. Теперь он не активен"
	GIVE_SUCCESSFULLY_DELETE_message    = "Конкурс успешно удален, он больше не будет выводится в списке твоих конкурсов"
	SELECT_PROPERTY_TO_EDIT_message     = "Что ты хочешь изменить?"

	CANNOT_GET_ADMINS_message                 = "Не могу получить список админов. Обратитесь к разработчику"
	CANNOT_CREATE_USER_message                = "Не могу создать пользователя. Обратитесь к разработчику"
	CANNOT_FIND_USER_message                  = "Не могу найти пользователя в базе. Обратитесь к разработчику"
	CANNOT_SET_USER_state_message             = "Не могу установить состояние пользователя. Обратитесь к разработчику"
	CANNOT_GET_USER_state_message             = "Не могу получить состояние пользователя. Обратитесь к разработчику"
	CANNOT_GET_STATE_DATA_message             = "He могу получить данные из состояния. Обратитесь к разработчику"
	CANNOT_GET_USER_GIVES_message             = "Не могу получить конкурсы пользователя. Обратитесь к разработчику"
	CANNOT_CREATE_GIVE_message                = "Не могу создать конкурс. Обратитесь к разработчику"
	CANNOT_UPDATE_GIVE_message                = "Не могу обновить конкурс. Обратитесь к разработчику"
	CANNON_UPDATE_GIVE_ON_PUBLICATION_message = "Не могу обновить конкурс giveId=%d во время публикации. Обратитесь к разработчику"
	CANNOT_GET_GIVE_message                   = "Не могу получить данные конкурса. Обратитесь к разработчику"
	CANNOT_DOWNLOAD_IMAGE_message             = "Не могу скачать картинку. Обратитесь к разработчику"
	CANNOT_PARSE_TIME_message                 = "Не могу распознать дату %s. Проверь корректность введенных данных, дата в формате ДД.ММ.ГГГГ ЧЧ:ММ"
	CANNOT_PARSE_WINNERS_COUNT_message        = "Не могу распознать количество победителей, проверь корректность ввода. Дольжно быть целое положительное число, ,больше 0"
	CANNOT_PARSE_CHANNEL_message              = "Не могу распознать канал %d"
	CANNOT_CHECK_BOT_IS_ADMIN_message         = "Не могу проверить канал %d на админа"
	CANNOT_SEND_message                       = "Не могу отправить сообщение пользователю %s"
	CANNOT_PUBLISH_GIVE_message               = "Не погу опубликовать конкурс giveId=%d. Обратитесь к разработчику"
	FINISH_DATE_HAS_PASSED_message            = "Дата окончания конкурса прошла"
	FINISH_DATE_BEFORE_START_message          = "Дата окончания конкурса раньше старта"
	BOT_MUST_BE_ADMIN_message                 = "Бот должен быть админов в канале %d"
	I_DONT_UNDERSTAND_message                 = "Я тебя не понимаю"
	GIVE_FIELDS_MUST_BE_FILLED_message        = "Поля конкурса должны быть заполнены:\n%s"
)
