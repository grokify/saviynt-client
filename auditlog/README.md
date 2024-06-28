# Saviynt Audit Log

## Audit Log Retrieval

A popular use case for the Saviynt API is to retrieve audit log information.

The following is an example of a Runtime Analytics SQL query that can be used to retrieve Audit Log details.

```
SELECT
	ua.LOGINKEY,
	l.LOGINTIME,
	l.LOGOUTDATE,
	l.COMMENTS AS LOGIN_COMMENTS,
	ua.TYPEOFACCESS AS OBJECTTYPE,
	ua.OBJECTKEY AS OBJECTNAME,
	ua.ActionType AS ACTION,
	u.username AS ACCESSBY,
	ua.IPADDRESS,
	ua.OBJECT_ATTRIBUTE_NAME AS ATTRIBUTE,
	ua.OLD_VALUE AS OLDVALUE,
	ua.NEW_VALUE AS NEWVALUE,
	ua.EVENT_ID AS EVENTID,
	ua.DETAIL,
	ua.ACCESS_URL,
	ua.ACCESSTIME AS EVENT_TIME,
	ua.QUERY_PARAM
FROM
	users u,
	userlogin_access ua,
	userlogins l
WHERE
	l.loginkey = ua.LOGINKEY AND
	l.USERKEY = u.userkey AND
	ua.AccessTime >= (NOW() - INTERVAL ${timeFrame} Minute) AND
	ua.Detail is not NULL
```
