-- .header on
.mode col

-- sqlpeek
-- x=time y=reward c=value row=decay
select * from wire_pow_of_two
;

-- Interesting how you don't see a massive difference
-- in the plots or in the total count of unique powers
-- of two.

-- sqlpeek 
select distinct campaign_id, proglen, decay
from wire_pow_of_two
;

-- sqlpeek
-- x=logrew y=bucket row=campaign_id c=campaign_id
with
S as (
    select *, log(2, reward + pow(2,-11)) as logrew
    from wire_pow_of_two
),
A as (
    select *, ntile(500) over (order by logrew) bucket
    from S
    order by value
),
B as (
    select *, row_number() over () as rn 
    from A
    group by bucket
)
select * from B
;

-- sqlpeek
SELECT name FROM sqlite_master WHERE type='table';

-- sqlpeek
create view campaignparams as
select distinct decay, proglen, campaign_id from wire_pow_of_two ;

-- unique powers-of-two by campaign
-- sqlpeek  
with
A as ( select distinct campaign_id, log(2, value) lg2val  from wire_pow_of_two ),
B as ( select *, floor(lg2val)=lg2val isPo2 from A ),
C as ( select * from B where isPo2 = true),
campaignparams as (select distinct decay, proglen, campaign_id from wire_pow_of_two),
D as ( select C.*, P.decay, P.proglen from C join campaignparams P using (campaign_id)),
E as ( select campaign_id, decay, proglen, count() total_unique from D group by campaign_id )
-- select * from D -- unique powers of two for each campaign
select * from E -- count of D by campaign

-- where isPo2 = 1
-- group by campaign_id
-- order by cnt
;


-- select only rows from campaigns with more than 1k rows
create view myview as
with A as (
    select campaign_id, count() cnt 
    from wire_pow_of_two 
    group by campaign_id
),
B as (
    select *,
        row_number() over (partition by campaign_id order by time) as rid
    from wire_pow_of_two join A using (campaign_id) 
    where cnt > 1000
)
select * from B 
;
