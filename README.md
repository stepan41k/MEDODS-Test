Для запуска приложения прописывается: sudo ./build.sh\
Обращение к эндпоинтам по адресу: http://localhost:8020\


Время жизни Access токена 15 минут, Refresh токена - 30 дней.\
Payload Access токена хранит GUID, key и дату истечения токена.\
Payload Refresh токена хранит IP, key и дату истечения токена.


Description of the endpoints:

"/auth/new/{guid}" - создаёт пару Access и Refresh токенов, принимая в параметре запроса GUID пользователя, если пользователя не сущуствует, возвращается соответсвующая ошибка. Refresh Token хранится в БД в виде хэша токена, а Access Token устанавливается как HTTP-Only cookie. 

"/auth/refresh" - (выполнятеся если установлена cookie с access token-ом, если не установлена, требуется вызов для создания токенов) достаёт из access токена GUID, с помощью которого обращается в БД за Refresh токеном, затем выполняет проверки(валидность, ip, и проверка на парность Access и Refresh). 


Auxiliary endpoints (for manual testing): 

"/auth/logout" - снимает HTTP-Only cookie

"/user/new" - генерирует GUID и добавляет его в БД


Также реализовано 5 функциональных тестов