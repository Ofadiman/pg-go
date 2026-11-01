insert into users (id, created_at, updated_at, address, birthday, email, email_verified, first_name, gender, last_name)
select
  format('00000000-0000-0000-0000-%s', lpad(g::text, 12, '0'))::uuid,
  t.created_at,
  case when g % 4 = 0 then null else t.created_at + ((g * 11) % 365 + 1) * interval '1 day' end,
  jsonb_build_object(
    'city', (array['Omaha', 'Boston', 'Toronto', 'Berlin', 'Warsaw', 'Krakow', 'Lyon', 'Madrid', 'Porto', 'Oslo'])[1 + (g * 3) % 10],
    'country', (array['United States', 'United States', 'Canada', 'Germany', 'Poland', 'Poland', 'France', 'Spain', 'Portugal', 'Norway'])[1 + (g * 3) % 10],
    'street', (array['Main Street', 'Oak Avenue', 'Maple Road', 'Elm Street', 'Pine Lane', 'Cedar Drive', 'Lake View', 'Hill Road'])[1 + (g * 5) % 8],
    'street_number', (1 + (g * 13) % 200)::text,
    'zip_code', lpad(((g * 7919) % 100000)::text, 5, '0')
  ),
  timestamptz '1960-01-01 00:00:00+00'
    + (case when g % 2 = 1 then (g * 37) % 56 else (g * g) % 56 end) * interval '1 year'
    + (1 + (g * 53) % 364) * interval '1 day',
  format('user-%s@example.com', g),
  g % 3 <> 0,
  case
    when g % 2 = 0 then (array['Anna', 'Angela', 'Bella', 'Chloe', 'Emma', 'Grace', 'Isla', 'Katie', 'Mia', 'Olivia'])[1 + ((g / 2) * 3) % 10]
    else (array['Andrew', 'Anthony', 'Brian', 'Daniel', 'Felix', 'Henry', 'Jack', 'Liam', 'Noah', 'Peter'])[1 + (((g + 1) / 2) * 7) % 10]
  end,
  case when g % 2 = 0 then 'female' else 'male' end,
  (array['Adams', 'Baker', 'Carter', 'Davis', 'Evans', 'Foster', 'Garcia', 'Hughes', 'Irwin', 'Jones', 'King', 'Lopez', 'Miller', 'Nelson', 'Owens', 'Parker', 'Quinn', 'Reed', 'Smith', 'Turner', 'Walker', 'Young', 'Zimmer'])[1 + (g * 5) % 23]
from generate_series(1, 100) as g
  cross join lateral (
    select timestamptz '2015-01-01 00:00:00+00' + ((g * 389) % 1826) * interval '1 day' + ((g * 7) % 24) * interval '1 hour' as created_at
  ) as t;

insert into posts (id, user_id, created_at, updated_at, content, title)
select
  format('00000000-0000-0000-0001-%s', lpad(pn.p::text, 12, '0'))::uuid,
  u.id,
  d.created_at,
  case when (un.n + k) % 3 = 0 then null else d.created_at + ((un.n * k) % 30 + 1) * interval '1 day' end,
  format('Notes on %s from user %s, part %s.', lower(d.topic), un.n, k),
  format('%s #%s', d.topic, pn.p)
from users as u
  cross join lateral (select right(u.id::text, 12)::int as n) as un
  cross join lateral generate_series(1, un.n % 10) as k
  cross join lateral (select (un.n - 1) * 9 + k as p) as pn
  cross join lateral (
    select
      u.created_at + ((un.n * 31 + k * 97) % 365 + 1) * interval '1 day' + k * interval '1 hour' as created_at,
      (array['Go generics', 'Postgres indexes', 'Query builders', 'Testing with containers', 'Error handling', 'Context cancellation', 'JSON in Postgres', 'Transactions', 'Window functions', 'Recursive queries', 'Connection pooling', 'Migrations'])[1 + (pn.p * 5) % 12] as topic
  ) as d
where un.n <= 90;

insert into comments (id, user_id, post_id, parent_id, created_at, updated_at, content)
select
  format('00000000-0000-0000-0002-%s', lpad(c.num::text, 12, '0'))::uuid,
  format('00000000-0000-0000-0000-%s', lpad(((c.post_number * 7 + c.k * 13) % 90 + 1)::text, 12, '0'))::uuid,
  c.post_id,
  null,
  c.created_at,
  case when c.num % 4 = 0 then c.created_at + interval '1 day' end,
  (array['Great post!', 'Thanks for sharing.', 'I disagree with this.', 'Could you elaborate?', 'This helped me a lot.', 'Interesting take.', 'I have a question.', 'Nice explanation.', 'Not sure about this.', 'Bookmarked.'])[1 + c.num % 10]
from (
  select
    po.id as post_id,
    pn.post_number,
    k,
    row_number() over (order by pn.post_number, k) as num,
    po.created_at + ((pn.post_number * 13 + k * 7) % 30 + 1) * interval '1 day' + k * interval '1 hour' as created_at
  from posts as po
    cross join lateral (select right(po.id::text, 12)::int as post_number) as pn
    cross join lateral generate_series(1, pn.post_number % 5) as k
) as c;

insert into comments (id, user_id, post_id, parent_id, created_at, updated_at, content)
select
  format('00000000-0000-0000-0002-%s', lpad(r.num::text, 12, '0'))::uuid,
  format('00000000-0000-0000-0000-%s', lpad(((r.parent_number * 11 + r.k * 17) % 90 + 1)::text, 12, '0'))::uuid,
  r.post_id,
  r.parent_id,
  r.created_at,
  case when r.num % 4 = 0 then r.created_at + interval '1 day' end,
  (array['Great post!', 'Thanks for sharing.', 'I disagree with this.', 'Could you elaborate?', 'This helped me a lot.', 'Interesting take.', 'I have a question.', 'Nice explanation.', 'Not sure about this.', 'Bookmarked.'])[1 + r.num % 10]
from (
  select
    co.id as parent_id,
    co.post_id,
    cn.num as parent_number,
    k,
    (select count(*) from comments) + row_number() over (order by cn.num, k) as num,
    co.created_at + ((cn.num * 7 + k) % 5 + 1) * interval '1 hour' as created_at
  from comments as co
    cross join lateral (select right(co.id::text, 12)::bigint as num) as cn
    cross join lateral generate_series(1, 2) as k
  where co.parent_id is null and cn.num % 3 = 0
) as r;

insert into comments (id, user_id, post_id, parent_id, created_at, updated_at, content)
select
  format('00000000-0000-0000-0002-%s', lpad(r.num::text, 12, '0'))::uuid,
  format('00000000-0000-0000-0000-%s', lpad(((r.parent_number * 11 + r.k * 17) % 90 + 1)::text, 12, '0'))::uuid,
  r.post_id,
  r.parent_id,
  r.created_at,
  case when r.num % 4 = 0 then r.created_at + interval '1 day' end,
  (array['Great post!', 'Thanks for sharing.', 'I disagree with this.', 'Could you elaborate?', 'This helped me a lot.', 'Interesting take.', 'I have a question.', 'Nice explanation.', 'Not sure about this.', 'Bookmarked.'])[1 + r.num % 10]
from (
  select
    co.id as parent_id,
    co.post_id,
    cn.num as parent_number,
    k,
    (select count(*) from comments) + row_number() over (order by cn.num, k) as num,
    co.created_at + ((cn.num * 7 + k) % 5 + 1) * interval '1 hour' as created_at
  from comments as co
    cross join lateral (select right(co.id::text, 12)::bigint as num) as cn
    cross join lateral generate_series(1, 1) as k
  where co.parent_id is not null and cn.num % 2 = 0
) as r;

insert into comments (id, user_id, post_id, parent_id, created_at, updated_at, content)
select v.id::uuid, v.user_id::uuid, po.id, v.parent_id::uuid, po.created_at + v.after, null, v.content
from posts as po
  cross join (
    values
      ('00000000-0000-0000-0002-900000000001', '00000000-0000-0000-0000-000000000002', null, interval '1 day', 'Which query builder should I pick?'),
      ('00000000-0000-0000-0002-900000000002', '00000000-0000-0000-0000-000000000003', '00000000-0000-0000-0002-900000000001', interval '1 day 1 hour', 'Depends on how much SQL you want to write.'),
      ('00000000-0000-0000-0002-900000000003', '00000000-0000-0000-0000-000000000004', '00000000-0000-0000-0002-900000000001', interval '1 day 2 hours', 'I would start with plain SQL.'),
      ('00000000-0000-0000-0002-900000000004', '00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0002-900000000002', interval '1 day 3 hours', 'As little as possible.'),
      ('00000000-0000-0000-0002-900000000005', '00000000-0000-0000-0000-000000000005', '00000000-0000-0000-0002-900000000002', interval '1 day 4 hours', 'Then an ORM might fit.'),
      ('00000000-0000-0000-0002-900000000006', '00000000-0000-0000-0000-000000000003', '00000000-0000-0000-0002-900000000003', interval '1 day 5 hours', 'Plain SQL does not compose well.'),
      ('00000000-0000-0000-0002-900000000007', '00000000-0000-0000-0000-000000000004', '00000000-0000-0000-0002-900000000004', interval '1 day 6 hours', 'Fair enough.'),
      ('00000000-0000-0000-0002-900000000008', '00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0002-900000000007', interval '1 day 7 hours', 'Thanks everyone!')
  ) as v (id, user_id, parent_id, after, content)
where po.id = '00000000-0000-0000-0001-000000000001';

insert into reactions (id, user_id, comment_id, created_at, emoji)
select
  format('00000000-0000-0000-0003-%s', lpad((row_number() over (order by cn.num, k))::text, 12, '0'))::uuid,
  format('00000000-0000-0000-0000-%s', lpad(((cn.num * 11 + k * 17) % 90 + 1)::text, 12, '0'))::uuid,
  co.id,
  co.created_at + k * interval '1 hour',
  (array['face_with_tears_of_joy', 'red_heart', 'smiling_face_with_heart_eyes', 'thumbs_up', 'smiling_face_with_smiling_eyes', 'fire', 'loudly_crying_face', 'hundred_points', 'party_popper', 'clapping_hands'])[1 + (cn.num + k) % 10]
from comments as co
  cross join lateral (select right(co.id::text, 12)::bigint as num) as cn
  cross join lateral generate_series(1, (cn.num % 4)::int) as k;
