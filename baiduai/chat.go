package baiduai

import (
	"fmt"
	"time"

	"github.com/byebyebruce/aggrsdk/pkg/util"
)

// https://ai.baidu.com/ai-doc/UNIT/Lkipmh0tz
const chatURL = `https://aip.baidubce.com/rpc/2.0/unit/service/v3/chat?access_token=%v`

type ChatContent struct {
	TerminalID string `json:"terminal_id"`
	Query      string `json:"query"`
}

type ChatRequest struct {
	Version   string      `json:"version"`
	ServiceID string      `json:"service_id"`
	SessionID string      `json:"session_id"`
	LogID     string      `json:"log_id"`
	Request   ChatContent `json:"request"`
}

type ChatResponse struct {
	Result struct {
		Version string `json:"version"`
		Context struct {
			SYSPRESUMEDSKILLS []string `json:"SYS_PRESUMED_SKILLS"`
			SYSPRESUMEDHIST   []string `json:"SYS_PRESUMED_HIST"`
			SYSVARS           struct {
			} `json:"SYS_VARS"`
		} `json:"context"`
		Timestamp string `json:"timestamp"`
		ServiceID string `json:"service_id"`
		SessionID string `json:"session_id"`
		LogID     string `json:"log_id"`
		RefID     string `json:"ref_id"`
		Responses []struct {
			Status int    `json:"status"`
			Msg    string `json:"msg"`
			Origin string `json:"origin"`
			Schema struct {
				Intents []struct {
					Slots []struct {
						SlotName   string `json:"slot_name"`
						SlotValues []struct {
							Confidence     float64       `json:"confidence"`
							Begin          int           `json:"begin"`
							Length         int           `json:"length"`
							OriginalWord   string        `json:"original_word"`
							NormalizedWord string        `json:"normalized_word"`
							SessionOffset  int           `json:"session_offset"`
							MergeMethod    string        `json:"merge_method"`
							SubSlots       []interface{} `json:"sub_slots"`
						} `json:"slot_values"`
					} `json:"slots"`
					IntentName       string  `json:"intent_name"`
					IntentConfidence float64 `json:"intent_confidence"`
					SluInfo          struct {
						SluIntent string `json:"slu_intent"`
						MatchInfo struct {
							FromWho string `json:"from_who"`
						} `json:"match_info"`
					} `json:"slu_info"`
				} `json:"intents"`
			} `json:"schema"`
			Actions []struct {
				Confidence float64       `json:"confidence"`
				Say        string        `json:"say"`
				Type       string        `json:"type"`
				Options    []interface{} `json:"options"`
				ActionID   string        `json:"action_id"`
				Img        []interface{} `json:"img"`
			} `json:"actions"`
			RawQuery          string `json:"raw_query"`
			SentimentAnalysis struct {
				Label string  `json:"label"`
				Pval  float64 `json:"pval"`
			} `json:"sentiment_analysis,omitempty"`
			LexicalAnalysis []struct {
				Term      string   `json:"term"`
				Weight    float64  `json:"weight"`
				Type      string   `json:"type"`
				Etypes    []string `json:"etypes"`
				BasicWord []string `json:"basic_word"`
			} `json:"lexical_analysis,omitempty"`
			SlotHistory []struct {
				SlotName   string `json:"slot_name"`
				SlotValues []struct {
					Confidence     float64       `json:"confidence"`
					Begin          int           `json:"begin"`
					Length         int           `json:"length"`
					OriginalWord   string        `json:"original_word"`
					NormalizedWord string        `json:"normalized_word"`
					SessionOffset  int           `json:"session_offset"`
					MergeMethod    string        `json:"merge_method"`
					SubSlots       []interface{} `json:"sub_slots"`
				} `json:"slot_values"`
			} `json:"slot_history"`
		} `json:"responses"`
	} `json:"result"`
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

// Chat 聊天api
func (this *BaiduAI) Chat(service string, text string) (string, error) {

	req := ChatRequest{
		Version:   "3.0",
		ServiceID: service,
		SessionID: "",
		LogID:     time.Now().String(),
		Request:   ChatContent{TerminalID: "mypi", Query: text},
	}
	resp := ChatResponse{}

	url := fmt.Sprintf(chatURL, this.Credentials.Access_token)
	err := util.HTTPPostJSON(url, req, &resp)
	if err != nil {
		return "", err
	}
	if resp.ErrorCode != 0 {
		return "", fmt.Errorf(resp.ErrorMsg)
	}
	if len(resp.Result.Responses) == 0 {
		return "", fmt.Errorf("no response")
	}
	if len(resp.Result.Responses[0].Actions) == 0 {
		return "", fmt.Errorf("no actor")
	}

	return resp.Result.Responses[0].Actions[0].Say, nil
}
