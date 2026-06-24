package interfaces

import (
	"deep-map-server/internal/delivery/binding"
	"deep-map-server/internal/interfaces/request"
	"encoding/json"
	"net/http"
)

// ---------------------------- HTTP Handlers --------------------------------

/**
* @description: 获取可供翻译的源语言列表
* @param {http.ResponseWriter} w
* @param {http.Request} r
* @return {error}
 */
func (c *TranslationController) handleGetSourceLanguages(w http.ResponseWriter, r *http.Request) {
	resp, err := c.getSourceLanguages(r.Context())
	if err != nil {
		binding.WriteError(w, err)
		return
	}
	binding.WriteJSON(w, http.StatusOK, resp)
}

/**
* @description: 获取可供翻译的目标语言列表
* @param {http.ResponseWriter} w
* @param {http.Request} r
* @return {error}
 */
func (c *TranslationController) handleGetTargetLanguages(w http.ResponseWriter, r *http.Request) {
	resp, err := c.getTargetLanguages(r.Context())
	if err != nil {
		binding.WriteError(w, err)
		return
	}
	binding.WriteJSON(w, http.StatusOK, resp)
}

/**
* @description: 翻译文本
* @param {http.ResponseWriter} w
* @param {http.Request} r
* @return {error}
 */
func (c *TranslationController) handleTranslateText(w http.ResponseWriter, r *http.Request) {
	var req request.TranslateTextRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		binding.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if err := binding.Validate(req); err != nil {
		binding.WriteError(w, err)
		return
	}
	resp, err := c.translateText(r.Context(), req)
	if err != nil {
		binding.WriteError(w, err)
		return
	}
	binding.WriteJSON(w, http.StatusOK, resp)
}

/**
* @description: 获取可供改写的写作风格列表
* @param {http.ResponseWriter} w
* @param {http.Request} r
* @return {error}
 */
func (c *TranslationController) handleGetWritingStyles(w http.ResponseWriter, r *http.Request) {
	resp, err := c.getWritingStyles(r.Context())
	if err != nil {
		binding.WriteError(w, err)
		return
	}
	binding.WriteJSON(w, http.StatusOK, resp)
}

/**
* @description: 获取可供改写的写作语气列表
* @param {http.ResponseWriter} w
* @param {http.Request} r
* @return {error}
 */
func (c *TranslationController) handleGetWritingTones(w http.ResponseWriter, r *http.Request) {
	resp, err := c.getWritingTones(r.Context())
	if err != nil {
		binding.WriteError(w, err)
		return
	}
	binding.WriteJSON(w, http.StatusOK, resp)
}

/**
* @description: 改写文本
* @param {http.ResponseWriter} w
* @param {http.Request} r
* @return {error}
 */
func (c *TranslationController) handleRephraseText(w http.ResponseWriter, r *http.Request) {
	var req request.RephraseTextRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		binding.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if err := binding.Validate(req); err != nil {
		binding.WriteError(w, err)
		return
	}
	resp, err := c.rephraseText(r.Context(), req)
	if err != nil {
		binding.WriteError(w, err)
		return
	}
	binding.WriteJSON(w, http.StatusOK, resp)
}

/**
* @description: 翻译文档
* @param {http.ResponseWriter} w
* @param {http.Request} r
* @return {error}
 */
func (c *TranslationController) handleTranslateDocument(w http.ResponseWriter, r *http.Request) {
	var req request.TranslateDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		binding.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if err := binding.Validate(req); err != nil {
		binding.WriteError(w, err)
		return
	}
	resp, err := c.translateDocument(r.Context(), req)
	if err != nil {
		binding.WriteError(w, err)
		return
	}
	binding.WriteJSON(w, http.StatusOK, resp)
}
