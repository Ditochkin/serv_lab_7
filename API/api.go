package API

import (
	"db_lab7/config"
	"db_lab7/db"
	"db_lab7/types"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type API struct {
	config *config.Config
	router *mux.Router
	store  *db.Store
}

func InitApi() (*API, error) {
	res := new(API)
	var err error
	res.config, err = config.GetConfig()

	if err != nil {
		return nil, err
	}

	res.router = mux.NewRouter()
	return res, nil
}

func (a *API) Start() error {
	a.configureRouter()
	a.configureDB()
	fmt.Println(a.store.Open())

	return http.ListenAndServe(a.config.Port, a.router)
}

func (a *API) Stop() {
	a.store.Close()
}

func (a *API) configureDB() {
	a.store = db.New(a.config)
}

func (a *API) configureRouter() {
	a.router.HandleFunc("/test", a.handleTest())

	a.router.HandleFunc("/add_publisher", a.handleAddPublisher())
	a.router.HandleFunc("/delete_publisher", a.handleDeletePublisher())
	a.router.HandleFunc("/change_publisher", a.handleChangePublisher())

	a.router.HandleFunc("/add_game_publisher", a.handleAddGamePublisher())
	a.router.HandleFunc("/delete_game_publisher", a.handleDeleteGamePublisher())

	a.router.HandleFunc("/delete_platform_by_year", a.handleDeletePlatformByYear())

	a.router.HandleFunc("/create_user", a.handleCreateUser())
	a.router.HandleFunc("/sign_in", a.handleSignIn())
	a.router.HandleFunc("/sign_out", a.handleSignOut())
}

func (a *API) handleSignOut() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		c := &http.Cookie{
			Name:     "session_token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		}
		http.SetCookie(writer, c)

		writer.WriteHeader(http.StatusOK)
	}
}

func (a *API) handleSignIn() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		_, _, err := a.GetIDAndRoleFromToken(writer, request)
		if err == nil {
			writer.WriteHeader(http.StatusOK)
			return
		}

		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "can't read body", http.StatusBadRequest)
			return
		}
		err = request.Body.Close()
		if err != nil {
			http.Error(writer, "can't close body", http.StatusInternalServerError)
			return
		}
		var usr types.User
		err = json.Unmarshal(body, &usr)
		if err != nil {
			http.Error(writer, "can't close body", http.StatusInternalServerError)
			return
		}
		token, err := a.generateTokensByCred(usr.Username, usr.Password)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		setTokenCookies(writer, token)
		writer.WriteHeader(http.StatusOK)
	}
}

func (a *API) handleCreateUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		_, role, err := a.GetIDAndRoleFromToken(writer, request)
		if err != nil {
			http.Error(writer, "You are not logged in. Sign In please", http.StatusBadRequest)
			return
		}
		if role != "admin" {
			http.Error(writer, "You are not admin and you have no right for this act.", http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "can't read body", http.StatusBadRequest)
			return
		}
		err = request.Body.Close()
		if err != nil {
			http.Error(writer, "can't close body", http.StatusInternalServerError)
			return
		}
		var usr types.User
		err = json.Unmarshal(body, &usr)
		if err != nil {
			http.Error(writer, "can't close body", http.StatusInternalServerError)
			return
		}
		_, err = a.store.Exec(db.CreateUserQuery, usr.Name, usr.Username, generatePasswordHash(usr.Password), usr.Role)
		if err != nil {
			if err.Error() == "UNIQUE constraint failed: users.Username" {
				http.Error(writer, "Username is already in use. Try to use another one.", http.StatusBadGateway)
				return
			}
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}

func (a *API) handleTest() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		a.Test()
	}
}

func (a *API) Test() error {
	rows, err := a.store.Query(db.SelectAllCountries)
	if err != nil {
		return err
	}
	defer rows.Close()
	var id int
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(id)
	}
	return nil
}

func (a *API) handleAddPublisher() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		_, _, err := a.GetIDAndRoleFromToken(writer, request)
		if err != nil {
			http.Error(writer, "You are not logged in. Sign In please", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "error in reading request", http.StatusBadRequest)
			return
		}

		err = request.Body.Close()
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		var publisher types.Publisher
		err = json.Unmarshal(body, &publisher)
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		if publisher.PublisherName == "" {
			http.Error(writer, "publisherName is empty", http.StatusInternalServerError)
			return
		}

		_, err = a.store.Exec(db.AddPublisherQuery, publisher.PublisherName)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}

func (a *API) handleDeletePublisher() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		_, _, err := a.GetIDAndRoleFromToken(writer, request)
		if err != nil {
			http.Error(writer, "You are not logged in. Sign In please", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(request.Body)

		if err != nil {
			http.Error(writer, "error in reading request", http.StatusBadRequest)
			return
		}

		err = request.Body.Close()
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		var publisher types.Publisher
		err = json.Unmarshal(body, &publisher)
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}
		if publisher.Id < 0 {
			http.Error(writer, "publisherName is empty", http.StatusInternalServerError)
			return
		}

		_, err = a.store.Exec(db.DeletePublisherQuery, publisher.PublisherName, publisher.PublisherName)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}

func (a *API) handleChangePublisher() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		_, _, err := a.GetIDAndRoleFromToken(writer, request)
		if err != nil {
			http.Error(writer, "You are not logged in. Sign In please", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "error in reading request", http.StatusBadRequest)
			return
		}

		err = request.Body.Close()
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		var publisher types.ChangePublisher
		err = json.Unmarshal(body, &publisher)
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		if publisher.PublisherName == "" || publisher.NewPublisherName == "" {
			http.Error(writer, "publisherName is empty", http.StatusInternalServerError)
			return
		}

		_, err = a.store.Exec(db.ChangePublisherQuery, publisher.NewPublisherName, publisher.PublisherName)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}

func (a *API) handleAddGamePublisher() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		_, _, err := a.GetIDAndRoleFromToken(writer, request)
		if err != nil {
			http.Error(writer, "You are not logged in. Sign In please", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "error in reading request", http.StatusBadRequest)
			return
		}

		err = request.Body.Close()
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		var publisher types.GamePublisher
		err = json.Unmarshal(body, &publisher)
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		if publisher.PublisherName == "" {
			http.Error(writer, "publisherName is empty", http.StatusInternalServerError)
			return
		}

		if publisher.GameName == "" {
			http.Error(writer, "gameName is empty", http.StatusInternalServerError)
			return
		}

		_, err = a.store.Exec(db.AddGamePublisherQuery, publisher.GameName, publisher.PublisherName)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}

func (a *API) handleDeleteGamePublisher() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		_, _, err := a.GetIDAndRoleFromToken(writer, request)
		if err != nil {
			http.Error(writer, "You are not logged in. Sign In please", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(request.Body)

		if err != nil {
			http.Error(writer, "error in reading request", http.StatusBadRequest)
			return
		}

		err = request.Body.Close()
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		var publisher types.GamePublisher
		err = json.Unmarshal(body, &publisher)
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		if publisher.PublisherName == "" {
			http.Error(writer, "publisherName is empty", http.StatusInternalServerError)
			return
		}

		if publisher.GameName == "" {
			http.Error(writer, "gameName is empty", http.StatusInternalServerError)
			return
		}

		_, err = a.store.Exec(db.DeleteGamePublisherQuery, publisher.GameName, publisher.PublisherName, publisher.GameName, publisher.PublisherName, publisher.GameName, publisher.PublisherName)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}

func (a *API) handleDeletePlatformByYear() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		_, _, err := a.GetIDAndRoleFromToken(writer, request)
		if err != nil {
			http.Error(writer, "You are not logged in. Sign In please", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "error in reading request", http.StatusBadRequest)
			return
		}

		err = request.Body.Close()
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		var publisher types.PlatformYear
		err = json.Unmarshal(body, &publisher)
		if err != nil {
			http.Error(writer, "wrong json body part", http.StatusInternalServerError)
			return
		}

		if publisher.Year <= 0 {
			http.Error(writer, "year is not positive", http.StatusInternalServerError)
			return
		}

		_, err = a.store.Exec(db.DeleteGamePlatformByYearQuery, publisher.Year, publisher.Year, publisher.Year)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}
