// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package model

import (
	"encoding/json"
	"io"
	"strconv"
)

const (
	WebsocketEventTyping                              = "typing"
	WebsocketEventPosted                              = "posted"
	WebsocketEventPostEdited                          = "post_edited"
	WebsocketEventPostDeleted                         = "post_deleted"
	WebsocketEventPostUnread                          = "post_unread"
	WebsocketEventChannelConverted                    = "channel_converted"
	WebsocketEventChannelCreated                      = "channel_created"
	WebsocketEventChannelDeleted                      = "channel_deleted"
	WebsocketEventChannelRestored                     = "channel_restored"
	WebsocketEventChannelUpdated                      = "channel_updated"
	WebsocketEventChannelMemberUpdated                = "channel_member_updated"
	WebsocketEventChannelSchemeUpdated                = "channel_scheme_updated"
	WebsocketEventDirectAdded                         = "direct_added"
	WebsocketEventGroupAdded                          = "group_added"
	WebsocketEventNewUser                             = "new_user"
	WebsocketEventAddedToTeam                         = "added_to_team"
	WebsocketEventLeaveTeam                           = "leave_team"
	WebsocketEventUpdateTeam                          = "update_team"
	WebsocketEventDeleteTeam                          = "delete_team"
	WebsocketEventRestoreTeam                         = "restore_team"
	WebsocketEventUpdateTeamScheme                    = "update_team_scheme"
	WebsocketEventUserAdded                           = "user_added"
	WebsocketEventUserUpdated                         = "user_updated"
	WebsocketEventUserRoleUpdated                     = "user_role_updated"
	WebsocketEventMemberroleUpdated                   = "memberrole_updated"
	WebsocketEventUserRemoved                         = "user_removed"
	WebsocketEventPreferenceChanged                   = "preference_changed"
	WebsocketEventPreferencesChanged                  = "preferences_changed"
	WebsocketEventPreferencesDeleted                  = "preferences_deleted"
	WebsocketEventEphemeralMessage                    = "ephemeral_message"
	WebsocketEventStatusChange                        = "status_change"
	WebsocketEventHello                               = "hello"
	WebsocketAuthenticationChallenge                  = "authentication_challenge"
	WebsocketEventReactionAdded                       = "reaction_added"
	WebsocketEventReactionRemoved                     = "reaction_removed"
	WebsocketEventResponse                            = "response"
	WebsocketEventEmojiAdded                          = "emoji_added"
	WebsocketEventChannelViewed                       = "channel_viewed"
	WebsocketEventMultipleChannelsViewed              = "multiple_channels_viewed"
	WebsocketEventPluginStatusesChanged               = "plugin_statuses_changed"
	WebsocketEventPluginEnabled                       = "plugin_enabled"
	WebsocketEventPluginDisabled                      = "plugin_disabled"
	WebsocketEventRoleUpdated                         = "role_updated"
	WebsocketEventLicenseChanged                      = "license_changed"
	WebsocketEventConfigChanged                       = "config_changed"
	WebsocketEventOpenDialog                          = "open_dialog"
	WebsocketEventGuestsDeactivated                   = "guests_deactivated"
	WebsocketEventUserActivationStatusChange          = "user_activation_status_change"
	WebsocketEventReceivedGroup                       = "received_group"
	WebsocketEventReceivedGroupAssociatedToTeam       = "received_group_associated_to_team"
	WebsocketEventReceivedGroupNotAssociatedToTeam    = "received_group_not_associated_to_team"
	WebsocketEventReceivedGroupAssociatedToChannel    = "received_group_associated_to_channel"
	WebsocketEventReceivedGroupNotAssociatedToChannel = "received_group_not_associated_to_channel"
	WebsocketEventGroupMemberDelete                   = "group_member_deleted"
	WebsocketEventGroupMemberAdd                      = "group_member_add"
	WebsocketEventSidebarCategoryCreated              = "sidebar_category_created"
	WebsocketEventSidebarCategoryUpdated              = "sidebar_category_updated"
	WebsocketEventSidebarCategoryDeleted              = "sidebar_category_deleted"
	WebsocketEventSidebarCategoryOrderUpdated         = "sidebar_category_order_updated"
	WebsocketWarnMetricStatusReceived                 = "warn_metric_status_received"
	WebsocketWarnMetricStatusRemoved                  = "warn_metric_status_removed"
	WebsocketEventCloudPaymentStatusUpdated           = "cloud_payment_status_updated"
	WebsocketEventCloudSubscriptionChanged            = "cloud_subscription_changed"
	WebsocketEventThreadUpdated                       = "thread_updated"
	WebsocketEventThreadFollowChanged                 = "thread_follow_changed"
	WebsocketEventThreadReadChanged                   = "thread_read_changed"
	WebsocketFirstAdminVisitMarketplaceStatusReceived = "first_admin_visit_marketplace_status_received"
	WebsocketEventDraftCreated                        = "draft_created"
	WebsocketEventDraftUpdated                        = "draft_updated"
	WebsocketEventDraftDeleted                        = "draft_deleted"
	WebsocketEventAcknowledgementAdded                = "post_acknowledgement_added"
	WebsocketEventAcknowledgementRemoved              = "post_acknowledgement_removed"
	WebsocketEventPersistentNotificationTriggered     = "persistent_notification_triggered"
	WebsocketEventHostedCustomerSignupProgressUpdated = "hosted_customer_signup_progress_updated"
)

type WebSocketMessage interface {
	ToJSON() ([]byte, error)
	IsValid() bool
	EventType() string
}

type WebsocketBroadcast struct {
	OmitUsers             map[string]bool `json:"omit_users"`                        // broadcast is omitted for users listed here
	UserId                string          `json:"user_id"`                           // broadcast only occurs for this user
	ChannelId             string          `json:"channel_id"`                        // broadcast only occurs for users in this channel
	TeamId                string          `json:"team_id"`                           // broadcast only occurs for users in this team
	ConnectionId          string          `json:"connection_id"`                     // broadcast only occurs for this connection
	OmitConnectionId      string          `json:"omit_connection_id"`                // broadcast is omitted for this connection
	ContainsSanitizedData bool            `json:"contains_sanitized_data,omitempty"` // broadcast only occurs for non-sysadmins
	ContainsSensitiveData bool            `json:"contains_sensitive_data,omitempty"` // broadcast only occurs for sysadmins
	// ReliableClusterSend indicates whether or not the message should
	// be sent through the cluster using the reliable, TCP backed channel.
	ReliableClusterSend bool `json:"-"`

	// BroadcastHooks is a slice of hooks IDs used to process events before sending them on individual connections. The
	// IDs should be understood by the WebSocket code.
	//
	// This field should never be sent to the client.
	BroadcastHooks []string `json:"broadcast_hooks,omitempty"`
	// BroadcastHookArgs is a slice of named arguments for each hook invocation. The index of each entry corresponds to
	// the index of a hook ID in BroadcastHooks
	//
	// This field should never be sent to the client.
	BroadcastHookArgs []map[string]any `json:"broadcast_hook_args,omitempty"`
}

func (wb *WebsocketBroadcast) copy() *WebsocketBroadcast {
	if wb == nil {
		return nil
	}

	var c WebsocketBroadcast
	if wb.OmitUsers != nil {
		c.OmitUsers = make(map[string]bool, len(wb.OmitUsers))
		for k, v := range wb.OmitUsers {
			c.OmitUsers[k] = v
		}
	}
	c.UserId = wb.UserId
	c.ChannelId = wb.ChannelId
	c.TeamId = wb.TeamId
	c.OmitConnectionId = wb.OmitConnectionId
	c.ContainsSanitizedData = wb.ContainsSanitizedData
	c.ContainsSensitiveData = wb.ContainsSensitiveData
	c.BroadcastHooks = wb.BroadcastHooks
	c.BroadcastHookArgs = wb.BroadcastHookArgs

	return &c
}

func (wb *WebsocketBroadcast) AddHook(hookID string, hookArgs map[string]any) {
	wb.BroadcastHooks = append(wb.BroadcastHooks, hookID)
	wb.BroadcastHookArgs = append(wb.BroadcastHookArgs, hookArgs)
}

type precomputedWebSocketEventJSON struct {
	Event     json.RawMessage
	Data      json.RawMessage
	Broadcast json.RawMessage
}

func (p *precomputedWebSocketEventJSON) copy() *precomputedWebSocketEventJSON {
	if p == nil {
		return nil
	}

	var c precomputedWebSocketEventJSON

	if p.Event != nil {
		c.Event = make([]byte, len(p.Event))
		copy(c.Event, p.Event)
	}

	if p.Data != nil {
		c.Data = make([]byte, len(p.Data))
		copy(c.Data, p.Data)
	}

	if p.Broadcast != nil {
		c.Broadcast = make([]byte, len(p.Broadcast))
		copy(c.Broadcast, p.Broadcast)
	}

	return &c
}

// webSocketEventJSON mirrors WebSocketEvent to make some of its unexported fields serializable
type webSocketEventJSON struct {
	Event     string              `json:"event"`
	Data      map[string]any      `json:"data"`
	Broadcast *WebsocketBroadcast `json:"broadcast"`
	Sequence  int64               `json:"seq"`
}

type WebSocketEvent struct {
	event           string
	data            map[string]any
	broadcast       *WebsocketBroadcast
	sequence        int64
	precomputedJSON *precomputedWebSocketEventJSON
}

// PrecomputeJSON precomputes and stores the serialized JSON for all fields other than Sequence.
// This makes ToJSON much more efficient when sending the same event to multiple connections.
func (ev *WebSocketEvent) PrecomputeJSON() *WebSocketEvent {
	evCopy := ev.Copy()
	event, _ := json.Marshal(evCopy.event)
	data, _ := json.Marshal(evCopy.data)
	broadcast, _ := json.Marshal(evCopy.broadcast)
	evCopy.precomputedJSON = &precomputedWebSocketEventJSON{
		Event:     json.RawMessage(event),
		Data:      json.RawMessage(data),
		Broadcast: json.RawMessage(broadcast),
	}
	return evCopy
}

func (ev *WebSocketEvent) RemovePrecomputedJSON() *WebSocketEvent {
	evCopy := ev.DeepCopy()
	evCopy.precomputedJSON = nil
	return evCopy
}

// WithoutBroadcastHooks gets the broadcast hook information from a WebSocketEvent and returns the event without that.
// If the event has broadcast hooks, a copy of the event is returned. Otherwise, the original event is returned. This
// is intended to be called before the event is sent to the client.
func (ev *WebSocketEvent) WithoutBroadcastHooks() (*WebSocketEvent, []string, []map[string]any) {
	hooks := ev.broadcast.BroadcastHooks
	hookArgs := ev.broadcast.BroadcastHookArgs

	if len(hooks) == 0 && len(hookArgs) == 0 {
		return ev, hooks, hookArgs
	}

	evCopy := ev.Copy()
	evCopy.broadcast = ev.broadcast.copy()

	evCopy.broadcast.BroadcastHooks = nil
	evCopy.broadcast.BroadcastHookArgs = nil

	return evCopy, hooks, hookArgs
}

func (ev *WebSocketEvent) Add(key string, value any) {
	ev.data[key] = value
}

func NewWebSocketEvent(event, teamId, channelId, userId string, omitUsers map[string]bool, omitConnectionId string) *WebSocketEvent {
	return &WebSocketEvent{
		event: event,
		data:  make(map[string]any),
		broadcast: &WebsocketBroadcast{
			TeamId:           teamId,
			ChannelId:        channelId,
			UserId:           userId,
			OmitUsers:        omitUsers,
			OmitConnectionId: omitConnectionId},
	}
}

func (ev *WebSocketEvent) Copy() *WebSocketEvent {
	evCopy := &WebSocketEvent{
		event:           ev.event,
		data:            ev.data,
		broadcast:       ev.broadcast,
		sequence:        ev.sequence,
		precomputedJSON: ev.precomputedJSON,
	}
	return evCopy
}

func (ev *WebSocketEvent) DeepCopy() *WebSocketEvent {
	evCopy := &WebSocketEvent{
		event:           ev.event,
		data:            copyMap(ev.data),
		broadcast:       ev.broadcast.copy(),
		sequence:        ev.sequence,
		precomputedJSON: ev.precomputedJSON.copy(),
	}
	return evCopy
}

func copyMap[K comparable, V any](m map[K]V) map[K]V {
	dataCopy := make(map[K]V, len(m))
	for k, v := range m {
		dataCopy[k] = v
	}
	return dataCopy
}

func (ev *WebSocketEvent) GetData() map[string]any {
	return ev.data
}

func (ev *WebSocketEvent) GetBroadcast() *WebsocketBroadcast {
	return ev.broadcast
}

func (ev *WebSocketEvent) GetSequence() int64 {
	return ev.sequence
}

func (ev *WebSocketEvent) SetEvent(event string) *WebSocketEvent {
	evCopy := ev.Copy()
	evCopy.event = event
	return evCopy
}

func (ev *WebSocketEvent) SetData(data map[string]any) *WebSocketEvent {
	evCopy := ev.Copy()
	evCopy.data = data
	return evCopy
}

func (ev *WebSocketEvent) SetBroadcast(broadcast *WebsocketBroadcast) *WebSocketEvent {
	evCopy := ev.Copy()
	evCopy.broadcast = broadcast
	return evCopy
}

func (ev *WebSocketEvent) SetSequence(seq int64) *WebSocketEvent {
	evCopy := ev.Copy()
	evCopy.sequence = seq
	return evCopy
}

func (ev *WebSocketEvent) IsValid() bool {
	return ev.event != ""
}

func (ev *WebSocketEvent) EventType() string {
	return ev.event
}

func (ev *WebSocketEvent) ToJSON() ([]byte, error) {
	if ev.precomputedJSON != nil {
		return ev.precomputedJSONBuf(), nil
	}
	return json.Marshal(webSocketEventJSON{
		ev.event,
		ev.data,
		ev.broadcast,
		ev.sequence,
	})
}

// Encode encodes the event to the given encoder.
func (ev *WebSocketEvent) Encode(enc *json.Encoder, buf io.Writer) error {
	if ev.precomputedJSON != nil {
		_, err := buf.Write(ev.precomputedJSONBuf())
		return err
	}

	return enc.Encode(webSocketEventJSON{
		ev.event,
		ev.data,
		ev.broadcast,
		ev.sequence,
	})
}

// We write optimal code here sacrificing readability for
// performance.
func (ev *WebSocketEvent) precomputedJSONBuf() []byte {
	return []byte(`{"event": ` +
		string(ev.precomputedJSON.Event) +
		`, "data": ` +
		string(ev.precomputedJSON.Data) +
		`, "broadcast": ` +
		string(ev.precomputedJSON.Broadcast) +
		`, "seq": ` +
		strconv.Itoa(int(ev.sequence)) +
		`}`)
}

func WebSocketEventFromJSON(data io.Reader) (*WebSocketEvent, error) {
	var ev WebSocketEvent
	var o webSocketEventJSON
	if err := json.NewDecoder(data).Decode(&o); err != nil {
		return nil, err
	}
	ev.event = o.Event
	if u, ok := o.Data["user"]; ok {
		// We need to convert to and from JSON again
		// because the user is in the form of a map[string]any.
		buf, err := json.Marshal(u)
		if err != nil {
			return nil, err
		}

		var user User
		if err = json.Unmarshal(buf, &user); err != nil {
			return nil, err
		}
		o.Data["user"] = &user
	}
	ev.data = o.Data
	ev.broadcast = o.Broadcast
	ev.sequence = o.Sequence
	return &ev, nil
}

// WebSocketResponse represents a response received through the WebSocket
// for a request made to the server. This is available through the ResponseChannel
// channel in WebSocketClient.
type WebSocketResponse struct {
	Status   string         `json:"status"`              // The status of the response. For example: OK, FAIL.
	SeqReply int64          `json:"seq_reply,omitempty"` // A counter which is incremented for every response sent.
	Data     map[string]any `json:"data,omitempty"`      // The data contained in the response.
	Error    *AppError      `json:"error,omitempty"`     // A field that is set if any error has occurred.
}

func (m *WebSocketResponse) Add(key string, value any) {
	m.Data[key] = value
}

func NewWebSocketResponse(status string, seqReply int64, data map[string]any) *WebSocketResponse {
	return &WebSocketResponse{Status: status, SeqReply: seqReply, Data: data}
}

func NewWebSocketError(seqReply int64, err *AppError) *WebSocketResponse {
	return &WebSocketResponse{Status: StatusFail, SeqReply: seqReply, Error: err}
}

func (m *WebSocketResponse) IsValid() bool {
	return m.Status != ""
}

func (m *WebSocketResponse) EventType() string {
	return WebsocketEventResponse
}

func (m *WebSocketResponse) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func WebSocketResponseFromJSON(data io.Reader) (*WebSocketResponse, error) {
	var o *WebSocketResponse
	return o, json.NewDecoder(data).Decode(&o)
}
