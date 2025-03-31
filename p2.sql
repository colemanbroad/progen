-- .header on
.mode col

-- campaign params
select distinct decay, proglen, campaign_id from wire_pow_of_two ;

-- .exit

-- unique powers-of-two by campaign
with A as ( select distinct campaign_id, log(2, value) val  from wire_pow_of_two ),
B as ( select *, floor(val)=val isPo2 from A ),
C as ( select decay, proglen, B.* from wire_pow_of_two right join B using (campaign_id))

select campaign_id, proglen, decay, count() cnt 
from C
where isPo2 = 1
group by campaign_id
order by cnt
;

.exit

select * from myview limit 3;

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
