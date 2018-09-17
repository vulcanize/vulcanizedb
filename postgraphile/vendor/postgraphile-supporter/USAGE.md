# @graphile/plugin-supporter

Thank you for being a supporter of PostGraphile!

To use this plugin:

1. install it with `yarn add` or `npm install` using the
   `https://...:...@git.graphile.com/.../postgraphile-supporter.git` URL from
   https://store.graphile.com - you must keep this URL confidential

2. tell postgraphile about the plugin via 
   `postgraphile --plugins @graphile/plugin-supporter`

3. to enable subscriptions, turn on the `--simple-subscriptions` option:
   `postgraphile --plugins @graphile/plugin-supporter --simple-subscriptions`

### Topic prefix

All topics are automatically prefixed with 'postgraphile:' but you can
customise this with the `pgSubscriptionPrefix` setting.

### Subscription security

By default, any user may subscribe to any topic, whether logged in or not, and
they will remain subscribed until they close the connection themselves. This
can cause a number of security issues; so we give you a method to implement
security around subscriptions.

By specifying `--subscription-authorization-function app_private.validate_subscription` on the PostGraphile CLI (or using the
`subscriptionAuthorizationFunction` option) you can have PostGraphile call the
function `app_private.validate_subscription(text)` to ensure that the user is
allowed to subscribe to the relevant topic (note: the `topic` argument WILL
be sent including the 'postgraphile:' prefix). The function will take the
following form:

```sql
CREATE FUNCTION app_hidden.validate_subscription(topic text)
RETURNS TEXT AS $$
BEGIN
  IF ... THEN
    RETURN ...::text;
  ELSE
    RAISE EXCEPTION 'Subscription denied' USING errcode = '.....';
  END IF;
END;
$$ LANGUAGE plpgsql VOLATILE SECURITY DEFINER;
```

You must define this function with your custom security logic. The function
must accept one text argument and either return a non-null text value, or throw
an error. The text value returned is used to tell the system when to cancel
the subscription - if you don't need this functionality then you may return a
_static_ unique value, e.g. generate a random UUID (manually) and then return
this same UUID over and over from your function, e.g.:

```sql
CREATE FUNCTION app_hidden.validate_subscription(topic text)
RETURNS TEXT AS $$
BEGIN
  IF ... THEN
    RETURN '4A2D27CD-7E67-4585-9AD8-50AF17D98E0B'::text;
  ELSE
    RAISE EXCEPTION 'Subscription denied' USING errcode = '.....';
  END IF;
END;
$$ LANGUAGE plpgsql VOLATILE SECURITY DEFINER;
```

You might want to make the topic a combination of things, for example the
subject type and identifier - e.g. 'channel:123'. If you do this then your
function could determine which subject the user is attempting to subscribe to,
check the user has access to that subject, and finally return a PostgreSQL
topic that will be published to in the event the user is kicked from the
channel, e.g.  `'postgraphile:channel:kick:123:987'` (assuming '987' is the id
of the current user).
