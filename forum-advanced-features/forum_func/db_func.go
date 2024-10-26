package forum_func

import (
	"database/sql"
	"fmt"
	"strconv"
)

func ExecSQLScript(databaseName, sqlScript string) (sql.Result, error) {
	DB, _ := sql.Open("sqlite3", databaseName) // Open the created SQLite File
	value, err := DB.Exec(sqlScript)
	if err != nil {
		fmt.Println(err)
	}
	DB.Close() // Defer Closing the database
	return value, err
}

func SelectSQLScript(databaseName, sqlScript string) *sql.Rows {
	DB, _ := sql.Open("sqlite3", databaseName) // Open the created SQLite File
	// fmt.Println("Select script:", sqlScript)
	values, err := DB.Query(sqlScript)
	if err != nil {
		fmt.Println(err)
	}
	DB.Close() // Defer Closing the database
	return values
}

type ForumTopic struct {
	TopicID               int
	TopicName             string
	TopicText             string
	TopicImg              string
	UserName              string
	CreationDate          string
	Likes                 int
	Hates                 int
	MyLikes               int
	MyHates               int
	Categories            string
	ContentCnt            int
	Flag                  int
	LastContentEditorDate string
	Status                string
	StatusText            string
}

type ForumTopicContent struct {
	UserName     string
	TopicID      int
	ContentID    int
	ContentName  string
	ContentText  string
	CreationDate string
	Likes        int
	Hates        int
	Flag         int
}

type Topic struct {
	TopicID   int
	TopicName string
	TopicText string
	TopicImg  string
}

type MyLikeDislikeTopic struct {
	TopicID   int
	TopicName string
	TLike     int
	THate     int
	CLike     int
	CHate     int
}

type LikeDisLike struct {
	Likes     int
	DisLikes  int
	MyLike    int
	MyDisLike int
}

type MyComment struct {
	TopicID      int
	TopicName    string
	ContentID    int
	ContentName  string
	ContentText  string
	CreationDate string
	Likes        int
	Hates        int
}

type ChangesAfterLogin struct {
	TopicID   int
	TopicName string
	Type      string
}

func MainTopics(databaseName, user, topics string) []ForumTopic {
	var ForumTopics []ForumTopic
	sqlScript := `	select
						h.topic_id,
						h.topic_name,
						h.topic_text,
						h.topic_img,
						u.user_name,
						strftime('%d.%m.%Y %H:%M', h.creation_date) as creation_date, 
						emotion.likes,
						emotion.hates,
						IFNULL((select distinct 1 from topic_emotion where status='ACTIVE' and emotion = 'Like' and user_email = '` + user + `' and topic_id = h.topic_id and content_id is null),0) mylikes,
						IFNULL((select distinct 1 from topic_emotion where status='ACTIVE' and emotion = 'Dislike' and user_email = '` + user + `' and topic_id = h.topic_id and content_id is null),0) myhates,
						(select group_concat(category_name, ', ') from topic_categories where topic_id = h.topic_id and status = 'ACTIVE') as categories,
						content.ContentCnt,
						case h.user_email when '` + user + `' then 1 else 0 end as flag,
						strftime('%d.%m.%Y %H:%M', content.LastContentEditorDate) as LastContentEditorDate,
						h.status,
						h.status_text
					from topic_header h 
					join users u 
						on h.user_email = u.user_email 
					left join (	select
									topic_id,
									count(content_id) as ContentCnt,
									max(status_date) as LastContentEditorDate
								from topic_content
								where status = 'ACTIVE'
								group by topic_id) as content
						on content.topic_id = h.topic_id
					left join (	select 
										topic_id,
										count(case emotion when 'Like' then emotion_id end) likes,
										count(case emotion when 'Dislike' then emotion_id end) hates
								from topic_emotion
								where status = 'ACTIVE'
									and content_id is null
									group by topic_id) emotion
						on emotion.topic_id = h.topic_id
					where 
					('AllTopics' = '` + topics + `' and (h.status != 'DELETED' or (h.user_email = '` + user + `' and h.status != 'DELETED')))
						or ('MyTopics' = '` + topics + `' and h.status != 'DELETED' and h.user_email = '` + user + `' )
						or ('PendingTopics' = '` + topics + `' and h.status = 'PENDING')`
	//    ('AllTopics' = '` + topics + `' and (h.status = 'ACTIVE' or (h.user_email = '` + user + `' and h.status != 'DELETED')))

	value := SelectSQLScript(databaseName, sqlScript)
	var user_name, topic_name, topic_text, topic_img, creation_date, categories, edit_date, status, status_text sql.NullString
	var topic_id, content_cnt, likes, hates, mylikes, myhates, flag sql.NullInt16

	for value.Next() {
		value.Scan(&topic_id, &topic_name, &topic_text, &topic_img, &user_name, &creation_date, &likes, &hates, &mylikes, &myhates, &categories, &content_cnt, &flag, &edit_date, &status, &status_text)
		mainTopics := ForumTopic{
			TopicID:               int(topic_id.Int16),
			TopicName:             topic_name.String,
			TopicText:             topic_text.String,
			TopicImg:              topic_img.String,
			UserName:              user_name.String,
			CreationDate:          creation_date.String,
			Likes:                 int(likes.Int16),
			Hates:                 int(hates.Int16),
			MyLikes:               int(mylikes.Int16),
			MyHates:               int(myhates.Int16),
			Categories:            categories.String,
			ContentCnt:            int(content_cnt.Int16),
			Flag:                  int(flag.Int16),
			LastContentEditorDate: edit_date.String,
			Status:                status.String,
			StatusText:            status_text.String}
		ForumTopics = append(ForumTopics, mainTopics)
	}
	return ForumTopics
}

func UpdateTopic(databaseName, user, id, name, text, img, removeImg string, categories []string) {
	maybeDeleteImg(id, removeImg, img)
	var sqlScript string
	imgValue := img
	if img == "remove" {
		imgValue = "false"
	}
	sqlScript = `update topic_header set
					topic_name = '` + name + `',
					topic_text = '` + text + `',
					topic_img = '` + imgValue + `',
					status_date = CURRENT_TIMESTAMP,
					status = 'PENDING' 
					where topic_id = ` + id + `
						and user_email = '` + user + `'
						and status != 'DELETED'`
	ExecSQLScript(databaseName, sqlScript)
	sqlScript = `DELETE FROM topic_categories WHERE topic_id = ` + id + `;`
	ExecSQLScript(databaseName, sqlScript)
	for _, value := range categories {
		sqlScript = `insert into topic_categories (topic_id, category_name) values (` + id + `, '` + value + `')`
		ExecSQLScript(databaseName, sqlScript)
	}
}

func UpdateSubTopic(databaseName, user, id, name, text string) {
	sqlScript := `update topic_content set
					content_name = '` + name + `',
					content_text = '` + text + `',
					status_date = CURRENT_TIMESTAMP 
					where content_id = ` + id + `
						and user_email = '` + user + `'
						and status = 'ACTIVE'`
	ExecSQLScript(databaseName, sqlScript)
}

func InsertTopic(databaseName, user, name, text, img string, categories []string) {
	sqlScript := `insert into topic_header (user_email, topic_name, topic_text, topic_img) values('` + user + `','` + name + `','` + text + `','` + img + `')`
	ExecSQLScript(databaseName, sqlScript)
	topic_id := GetLastRecord(databaseName, "topic_header")
	for _, v := range categories {
		sqlScript = `insert into topic_categories (topic_id, category_name) values('` + strconv.Itoa(topic_id) + `','` + v + `')`
		ExecSQLScript(databaseName, sqlScript)
	}
}

func GetLastRecord(databaseName, table string) int {
	var last_id int
	sqlScript := `select seq from sqlite_sequence where name='` + table + `';`
	value := SelectSQLScript(databaseName, sqlScript)
	var last_record sql.NullInt16
	for value.Next() {
		value.Scan(&last_record)
		last_id = int(last_record.Int16)
	}
	return last_id
}

func InsertSubTopic(databaseName, user, topicid, name, text string) {
	sqlScript := `insert into topic_content (user_email, topic_id, content_name, content_text) values('` + user + `','` + topicid + `','` + name + `','` + text + `')`
	ExecSQLScript(databaseName, sqlScript)
}

func DeleteTopic(databaseName, user, id string) {
	deleteImg(id)
	sqlScript := `update topic_header set
					status = 'DELETED',
					status_date = CURRENT_TIMESTAMP
					where status != 'DELETED'
						and user_email = '` + user + `'
						and topic_id = ` + id
	ExecSQLScript(databaseName, sqlScript)
}

func DeleteSubTopic(databaseName, user, id string) {
	sqlScript := `update topic_content set
					status = 'DELETED',
					status_date = CURRENT_TIMESTAMP
					where status = 'ACTIVE'
						and user_email = '` + user + `'
						and content_id = ` + id
	ExecSQLScript(databaseName, sqlScript)
}

func SubTopics(databaseName, user, id string) []ForumTopicContent {
	var ForumTopics []ForumTopicContent
	sqlScript := `select 
							case when u.status = 'DELETED' then u.user_name || ' (Inactive)' else u.user_name end as user_name, 
							h.topic_id,
							c.content_id, 
							c.content_name, 
							c.content_text, 
							strftime('%d.%m.%Y %H:%M', c.creation_date) as creation_date,
							count(distinct case when e.emotion = 'Like' then e.emotion_id end) likes,
							count(distinct case when e.emotion = 'Dislike' then e.emotion_id end) hates,
							case c.user_email when '` + user + `' then 1 else 0 end as flag
					from topic_header h
					join topic_content c 
						on h.topic_id = c.topic_id
						and c.status = 'ACTIVE'
					join users u 
						on c.user_email = u.user_email 
					left join topic_emotion e
						on e.topic_id = ` + id + `
						and e.content_id = c.content_id
						and e.status = 'ACTIVE'
					where h.topic_id = ` + id + `
					group by u.status, u.user_name, h.topic_id, c.content_id, h.topic_name, h.topic_text, c.content_name, c.content_text, c.creation_date, c.user_email`
	value := SelectSQLScript(databaseName, sqlScript)
	var user_name, content_name, content_text, creation_date sql.NullString
	var topic_id, content_id, likes, hates, flag sql.NullInt16

	for value.Next() {
		value.Scan(&user_name, &topic_id, &content_id, &content_name, &content_text, &creation_date, &likes, &hates, &flag)
		subTopics := ForumTopicContent{
			UserName:     user_name.String,
			TopicID:      int(topic_id.Int16),
			ContentID:    int(content_id.Int16),
			ContentName:  content_name.String,
			ContentText:  content_text.String,
			CreationDate: creation_date.String,
			Likes:        int(likes.Int16),
			Hates:        int(hates.Int16),
			Flag:         int(flag.Int16)}
		ForumTopics = append(ForumTopics, subTopics)
	}
	return ForumTopics
}

func GetTopic(databaseName, id string) Topic {
	var getTopics Topic
	sqlScript := `select topic_id, topic_name, topic_text, topic_img from topic_header where topic_id = ` + id
	value := SelectSQLScript(databaseName, sqlScript)
	var topic_id sql.NullInt16
	var topic_name sql.NullString
	var topic_text sql.NullString
	var topic_img sql.NullString
	for value.Next() {
		value.Scan(&topic_id, &topic_name, &topic_text, &topic_img)
		getTopics = Topic{
			TopicID:   int(topic_id.Int16),
			TopicName: topic_name.String,
			TopicText: topic_text.String,
			TopicImg:  topic_img.String}
	}

	return getTopics
}

func TopicLikes(databaseName, mainId, LikeType, user string) LikeDisLike {
	var LikesDisLikes LikeDisLike
	sqlScriptUpdate := `update topic_emotion set
						emotion = '` + LikeType + `',
						emotion_text = '` + LikeType + `',
						status = case when emotion = '` + LikeType + `' then 
									  case when status = 'ACTIVE' then 'DELETED' else 'ACTIVE' end
									  else 'ACTIVE' END,
						status_date = CURRENT_TIMESTAMP
					where content_id is null and topic_id = ` + mainId + ` and user_email = '` + user + `'`
	ExecSQLScript(databaseName, sqlScriptUpdate)
	sqlScriptInsert := `insert into topic_emotion (topic_id, user_email, emotion, emotion_text) 
						select ` + mainId + `, '` + user + `', '` + LikeType + `', '` + LikeType + `'
						where (select topic_id from topic_emotion where content_id is null and topic_id = ` + mainId + ` and user_email = '` + user + `') is null`
	ExecSQLScript(databaseName, sqlScriptInsert)
	sqlScript := `select count(case when emotion = 'Like' then emotion_id end) as likes,
						 count(case when emotion = 'Dislike' then emotion_id end) as dislikes,
						 count(distinct case when emotion = 'Like' and user_email = '` + user + `' then 1 end) mylike,
						 count(distinct case when emotion = 'Dislike' and user_email = '` + user + `' then 1 end) mydislike
		 from topic_emotion where status = 'ACTIVE' and content_id is null and topic_id = ` + mainId
	value := SelectSQLScript(databaseName, sqlScript)
	var likes, dislikes, mylike, mydislike sql.NullInt16
	for value.Next() {
		value.Scan(&likes, &dislikes, &mylike, &mydislike)
		LikesDisLikes = LikeDisLike{Likes: int(likes.Int16), DisLikes: int(dislikes.Int16), MyLike: int(mylike.Int16), MyDisLike: int(mydislike.Int16)}
	}
	return LikesDisLikes
}

func SubTopicLikes(databaseName, mainId, subId, LikeType, user string) LikeDisLike {
	var LikesDisLikes LikeDisLike
	sqlScriptUpdate := `update topic_emotion set
						emotion = '` + LikeType + `',
						emotion_text = '` + LikeType + `',
						status = case when emotion = '` + LikeType + `' then 
									  case when status = 'ACTIVE' then 'DELETED' else 'ACTIVE' end
									  else 'ACTIVE' END,
						status_date = CURRENT_TIMESTAMP
					where content_id = ` + subId + ` and topic_id = ` + mainId + ` and user_email = '` + user + `'`
	ExecSQLScript(databaseName, sqlScriptUpdate)
	sqlScriptInsert := `insert into topic_emotion (topic_id, content_id, user_email, emotion, emotion_text) 
						select ` + mainId + `, ` + subId + `, '` + user + `', '` + LikeType + `', '` + LikeType + `'
						where (select topic_id from topic_emotion where content_id = ` + subId + ` and topic_id = ` + mainId + ` and user_email = '` + user + `') is null`
	ExecSQLScript(databaseName, sqlScriptInsert)

	sqlScript := `select count(distinct case when emotion = 'Like' then emotion_id end) as likes,
						 count(distinct case when emotion = 'Dislike' then emotion_id end) as dislikes,
 						 count(distinct case when emotion = 'Like' and user_email = '` + user + `' then 1 end) mylike,
						 count(distinct case when emotion = 'Dislike' and user_email = '` + user + `' then 1 end) mydislike
		 from topic_emotion where status = 'ACTIVE' and content_id = ` + subId + ` and topic_id = ` + mainId
	value := SelectSQLScript(databaseName, sqlScript)
	var likes, dislikes, mylike, mydislike sql.NullInt16
	for value.Next() {
		value.Scan(&likes, &dislikes, &mylike, &mydislike)
		LikesDisLikes = LikeDisLike{Likes: int(likes.Int16), DisLikes: int(dislikes.Int16), MyLike: int(mylike.Int16), MyDisLike: int(mydislike.Int16)}
	}
	return LikesDisLikes
}

func ReadSetupTable(databaseName, setup_type string) []string {
	var values []string
	sqlScript := `select setup_value from setup where setup_type = '` + setup_type + `'`
	value := SelectSQLScript(databaseName, sqlScript)
	var setup_value sql.NullString
	for value.Next() {
		value.Scan(&setup_value)
		values = append(values, setup_value.String)
	}
	return values
}

func ReadTopicUsers(databaseName, user string) []string {
	var values []string
	sqlScript := `select distinct u.user_name from users u
	join topic_header h
		on u.user_email = h.user_email
	where h.status != 'DELETED'`
	// where h.status = 'ACTIVE' or (h.user_email = '` + user + `' and h.status != 'DELETED')`
	value := SelectSQLScript(databaseName, sqlScript)
	var setup_value sql.NullString
	for value.Next() {
		value.Scan(&setup_value)
		values = append(values, setup_value.String)
	}
	return values
}

func MyLikesDislikes(databaseName, user string) []MyLikeDislikeTopic {
	var values []MyLikeDislikeTopic
	sqlScript := `select h.topic_id, h.topic_name,
 						 count(distinct case when emotion = 'Like' and e.content_id is null then e.emotion_id end) tlike,
						 count(distinct case when emotion = 'Dislike' and e.content_id is null then e.emotion_id end) tdislike,
 						 count(distinct case when emotion = 'Like' and e.content_id is not null then e.emotion_id end) clike,
						 count(distinct case when emotion = 'Dislike' and e.content_id is not null then e.emotion_id end) cdislike
	from topic_header h
	join topic_emotion e
		on h.topic_id = e.topic_id
	where h.status != 'DELETED' and e.user_email = '` + user + `' and e.status != 'DELETED'
	group by h.topic_id, h.topic_name
	order by 2`
	value := SelectSQLScript(databaseName, sqlScript)
	var topic_name sql.NullString
	var topic_id, tlike, thate, clike, chate sql.NullInt16
	for value.Next() {
		value.Scan(&topic_id, &topic_name, &tlike, &thate, &clike, &chate)
		LikeValues := MyLikeDislikeTopic{
			TopicID:   int(topic_id.Int16),
			TopicName: topic_name.String,
			TLike:     int(tlike.Int16),
			THate:     int(thate.Int16),
			CLike:     int(clike.Int16),
			CHate:     int(chate.Int16)}
		values = append(values, LikeValues)
	}
	return values
}

func MyComments(databaseName, user string) []MyComment {
	var values []MyComment
	sqlScript := `select distinct h.topic_id, h.topic_name, c.content_id, c.content_name, c.content_text, strftime('%d.%m.%Y %H:%M', c.creation_date) as creation_date, emotion.likes, emotion.hates
 	from topic_header h
	join topic_content c
		on h.topic_id = c.topic_id
	left join (	select 
					content_id,
					count(case emotion when 'Like' then emotion_id end) likes,
					count(case emotion when 'Dislike' then emotion_id end) hates
				from topic_emotion
				where status = 'ACTIVE'
					and content_id is not null
					group by content_id) emotion
		on emotion.content_id = c.content_id
	where h.status != 'DELETED' and c.user_email = '` + user + `' and c.status != 'DELETED'`
	value := SelectSQLScript(databaseName, sqlScript)
	var topic_name, content_name, content_text, creation_date sql.NullString
	var topic_id, content_id, likes, hates sql.NullInt16
	for value.Next() {
		value.Scan(&topic_id, &topic_name, &content_id, &content_name, &content_text, &creation_date, &likes, &hates)
		Comments := MyComment{
			TopicID:      int(topic_id.Int16),
			TopicName:    topic_name.String,
			ContentID:    int(content_id.Int16),
			ContentName:  content_name.String,
			ContentText:  content_text.String,
			CreationDate: creation_date.String,
			Likes:        int(likes.Int16),
			Hates:        int(hates.Int16),
		}
		values = append(values, Comments)
	}
	return values
}

func ChangesAfterLastLogin(databaseName, user string) []ChangesAfterLogin {
	var values []ChangesAfterLogin
	sqlScript := `select distinct h.topic_id, h.topic_name, 'New Comments on topic' as type
	from topic_header h
	join users u
		on u.user_email = h.user_email
	join topic_content c
		on h.topic_id = c.topic_id
		and u.login_time_prev <= c.status_date
	where h.status != 'DELETED' and u.user_email = '` + user + `' and c.status != 'DELETED'

	union all

	select distinct h.topic_id, h.topic_name, case when e.content_id is null then e.emotion ||' on topic' else e.emotion ||' on comments' end as type
	from topic_header h
	join users u
		on u.user_email = h.user_email
	join topic_emotion e
		on h.topic_id = e.topic_id
		and u.login_time_prev <= e.status_date
	where h.status != 'DELETED' and u.user_email = '` + user + `' and e.status != 'DELETED'	`
	value := SelectSQLScript(databaseName, sqlScript)
	var topic_name, change_type sql.NullString
	var topic_id sql.NullInt16
	for value.Next() {
		value.Scan(&topic_id, &topic_name, &change_type)
		Comments := ChangesAfterLogin{
			TopicID:   int(topic_id.Int16),
			TopicName: topic_name.String,
			Type:      change_type.String}
		values = append(values, Comments)
	}
	return values
}
