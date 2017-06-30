CREATE or replace FUNCTION check_permission(key uuid, object_type varchar(20), permiss varchar(2), out group_id integer, out username varchar(70), out domain_name varchar(40) )
 AS $$
DECLARE cnt integer;
BEGIN

	select
		 acl_tokens.group_id,
		 acl_tokens.domain,
		 acl_tokens.username
	 into group_id, domain_name, username
	from acl_tokens where acl_tokens.key = $1 limit 1;

	with recursive groups as (

	  --start with the "anchor" row
	  select
	    *
	  from acl_group
	  where	id = group_id

	  union all

	  select
	    acl_group.*
	  from acl_group
	  join groups on groups.parent_id = acl_group.id
	)
	select
	    count(p.*) INTO cnt
	from groups g
		inner join acl_permission p on g.Id = p.group_id
	where p.object_type = $2 and (p.perm && '{*}'::varchar(2)[] or p.perm @> ARRAY[permiss]::varchar(2)[]) ;

	IF cnt < 1 THEN
		group_id = null;
		domain_name = '';
		username = '';
	END IF;
END;
$$  LANGUAGE plpgsql
