INSERT INTO public.acl_group(
             name)
    VALUES ('root');

INSERT INTO public.acl_group(
             name)
    VALUES ('user');

INSERT INTO public.acl_group(
             name, parent_id)
    VALUES ('admin', (select id from acl_group where name = 'user'));

--

SELECT id , perm
FROM acl_permission;

  SELECT
  id,
  object_type,
  object_id,
  group_id,
  array_to_json(perm)
FROM acl_permission where group_id = (select acl_group.id from acl_group where acl_group.name = 'root')


delete from  acl_permission where perm is null

INSERT INTO public.acl_permission(
             object_type, group_id, perm)
    VALUES ('acl/resource', 1, '{*}');
INSERT INTO public.acl_permission(
             object_type, group_id, perm)
    VALUES ('acl/roles', 1, '{*}');

--

INSERT INTO public.acl_tokens(
             key, username, group_id)
    VALUES ('ec7ce853-d6ae-45cf-a479-f557fa9bfa9d', 'igor', 1);
