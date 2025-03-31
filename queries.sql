-- .headers on
.mode col

select *
from history_power_of_two
limit 15;

-- plot: reduce the 1.1M rows down to 50k
drop table history_binned ;
CREATE table history_binned as
select avg(log(2, value + 1e-8)) val, avg(reward) reward, avg(time) time, campaign_id, mut
from history_power_of_two
group by campaign_id, ceil(reward*100), ceil(time/10)
order by random()
limit 50000
;

-- Are there any unsuccessful campaigns?
with A as (select distinct campaign_id, reward as val from history_binned)
select campaign_id, count(val) count from A group by campaign_id
order by count 
;

-- subsample campaigns
drop table hb2;
create table hb2 as
select * from history_power_of_two
natural join (select distinct campaign_id from history_binned limit 5)
order by campaign_id
;

-- how many unique powers of two? 
with A as ( select distinct log(2, value) val from history_power_of_two ),
B as ( select val, floor(val)=val isPo2 from A )
select count() from B where isPo2 = 1
;

-- unique powers-of-two by campaign
with A as ( select distinct campaign_id, log(2, value) val, mut from history_power_of_two ),
B as ( select *, floor(val)=val isPo2 from A )
select campaign_id, mut, count() cnt from B where isPo2 = 1
group by campaign_id, mut
order by cnt
;

-- plot: cumulative dist of number of unique powers of two across a campaign
-- group by mutation type 
drop table distplot;
create table distplot as
with A as ( select distinct campaign_id, log(2, value) val, mut from history_power_of_two ),
B as ( select *, floor(val)=val isPo2 from A ),
C as (
    select campaign_id, mut, count() cnt from B where isPo2 = 1
    group by campaign_id, mut
    order by cnt
)
select cnt, mut || "-mutate" mut, row_number() over (partition by mut order by cnt) rowid
from C
;

