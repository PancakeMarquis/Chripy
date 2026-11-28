package main

import "net/http"

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {

	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to retrive chirps", err)
		return
	}
	var list_chirps []Chirp
	for _, chirp := range chirps {
		list_chirps = append(list_chirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID})
	}

	respondWithJSON(w, http.StatusOK, list_chirps)
}
