workspace {
    model {
        user = person "Пользователь"
        admin = person "Администратор"
        guest = person "Гость"

        mypills = softwareSystem "MyPills" {
            frontend = container "Пользовательский интерфейс"
            backend = container "Бэкенд" {
                uiLayer = group "Компонент реализации UI" {
                    authController = component "Контроллер аутентификации"
                    profileController = component "Контроллер профиля"
                    cabinetController = component "Контроллер управления аптечкой"
                    medicineController = component "Контроллер подбора лекарств"
                    adminController = component "Контроллер администратора"
                }

                blLayer = group "Компонент бизнес-логики" {
                    authService = component "Сервис аутентификации"
                    profileService = component "Сервис профиля"
                    medicineService = component "Сервис подбора лекарств"
                    cabinetService = component "Сервис управления аптечкой"
                    adminService = component "Сервис администратора"
                    notifyService = component "Сервис уведомлений"
                }

                daLayer = group "Компонент доступа к данным" {
                    userRepository = component "Репозиторий пользователей"
                    medicineRepository = component "Репозиторий лекарств"
                    cabinetRepository = component "Репозиторий аптечек"
                }
            }
            database = container "База данных" {
                tags "Database"
            }
            cache = container "Кэш"{
                tags "Database"
            }
        }

        user -> mypills "Использует для управления аптечкой и подбора лекарства по симптому"
        admin -> mypills "Администрирует справочник лекарств"
        guest -> mypills "Регистрируется в системе"

        user -> frontend "Смотрит наличие лекарств и ищет подходящие по симптому"
        admin -> frontend "Видит панель управления для редактирования справочника"
        guest -> frontend "Видит окно для решистрации или входа"
        frontend -> backend "Делает запросы к"
        backend -> database "Читает из и пишет в"
        backend -> cache "Читает из и пишет в"

        frontend -> authController "Вызывает"
        frontend -> profileController "Вызывает"
        frontend -> medicineController "Вызывает"
        frontend -> adminController "Вызывает"
        frontend -> cabinetController "Вызывает"
        authController -> authService "Использует"
        profileController -> profileService "Использует"
        cabinetController -> cabinetService "Использует"
        medicineController -> medicineService "Использует"
        adminController -> adminService "Использует"
        authService -> userRepository "Использует"
        profileService -> userRepository "Использует"
        medicineService -> medicineRepository "Использует"
        medicineService -> userRepository "Использует"
        medicineService -> cabinetRepository "Использует"
        cabinetService -> cabinetRepository "Использует"
        adminService -> medicineRepository "Использует"
        notifyService -> cabinetRepository "Использует"
        notifyService -> userRepository "Использует"
        userRepository -> database "SQL"
        medicineRepository -> database "SQL"
        cabinetRepository -> database "SQL"
    }

    views {
        systemContext mypills "Context" {
            include *
            autolayout lr
        }

        container mypills "Containers" {
            include *
            autolayout lr
        }

        component backend "Components" {
            include *
            autolayout lr
        }

        styles {
            element "Database" {
                shape cylinder
            }

            element "Integration" {
                stroke #999999
                strokeWidth 5
                opacity 70
            }
        }

        theme default
    }
}