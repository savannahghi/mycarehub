## Documentation

The authorization determines a request based on ``{subject, object, action}``, which means what ``subject`` can perform what ``action`` on what ``object``. Here the meanings are:

1. ``subject``: the logged-in user name
2. ``object``: the URL path for the web resource like "dataset1/item1"
3. ``action``: HTTP method like GET, POST, PUT, DELETE, or the high-level actions you defined like "read-file", "write-blog"

The request_definition is how we interact with Casbin. We define our authorization requests to be in the format we want e.g user_id, feature, action (with action being either edit or view

[request_definition]
r = user_id, feature, action

The policy_definition is how we store permissions. In Casbin’s terminology, “policy” is equivalent to our concept of a “permission.” We define our policies to be stored in the format user_id, feature, action, mirroring how we send requests to Casbin to verify permissions.
[policy_definition]
p = user_id, feature, action


The role_definition in this case is used to create the feature hierarchy. The first _ indicates the parent feature and the second _ indicates the child feature. This is a little opaque, but functionally this block defines a “grouping policy” in the format g(a, b) where we can basically store and verify the relationship between a pair of strings (which in our case are parent and child features names).
[role_definition]
g = _, _

The policy_effect defines what action is taken if a request is matched according to the matcher in the configuration(which we’ll take a look at later). In this block, we’re basically saying that if some policy is matched, the effect is allow.
[policy_effect]
e = some(where (p.eft == allow))

matchers defines the boolean expression needed to satisfy an authorization request with existing policies.