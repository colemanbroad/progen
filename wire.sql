.mode col

-- First, list all the tables.
select name FROM sqlite_master WHERE type='table';

-- controlled variables
select distinct prog_l from wiring ;
select distinct wr_decay from wiring ;
select distinct wr_nearby from wiring ;
select distinct n_prog from wiring ;
select distinct cheating from wiring ;

-- measured variables
select distinct depth from wiring order by depth ;
select distinct count from wiring order by count;

-- Sample the first 20
select * from wiring limit 20;

-- How many total rows?
select count() from wiring;

-- How many different `cheating` hyperparms did we try?
select distinct cheating from wiring limit 5;

-- Uneven counts?
select wr_decay, prog_l, sum(count) from wiring
group by wr_decay, prog_l
;

-- Uneven counts?
select wr_decay, prog_l, sum(count) from wiring
group by campaign_id
;

-- Uneven total nodes in wire groups! Need to rerun analyses.
select wr_decay, sum(count) from wiring
group by wr_decay
;

-- Let's study coarse grained distributions by wire decay, conditional
-- on program length.
select wr_decay, depth/10, sum(count) from wiring
where prog_l = 100
group by depth/10, wr_decay
order by wr_decay, depth/10 
limit 50
;

-- x=depth y=count c=wr_decay col=cheating share=Xy
select depth, count, wr_decay, cheating
from wiring
order by cheating
limit 5000
;

-- x=depth y=count c=prog_l col=cheating share=Xy
select depth, count, prog_l, cheating
from wiring
where wr_decay = 0.1
order by cheating
limit 5000
;

create table calaban as
select *, 'zero' as wr_decay from wiring where wr_nearby = 0 
union all
select *, wr_decay as wr_decay from wiring where wr_nearby = 1 
;

select * from calaban limit 10 ;

drop table calaban;
alter table calaban rename to calaban ;
alter table calaban drop column wr_decay;
alter table calaban rename column "wr_decay:1" to wr_decay;

-- select 1;

select distinct prog_l from calaban;

-- Let's compare the count of syms at different depths. We expect that
-- changeing the method of wiring will affect the depth distribution.
-- Notice how the x-axis is scaled!
-- x=depth y=count c=cheating row=prog_l

select distinct wr_decay from calaban;

select * from calaban
where "wr_decay" = 'zero'
limit 50
;

select count() from  wiring;




drop view counthist ;

create view counthist as
select depth, count as count, prog_l, wr_decay
    from wiring where wr_nearby = 1 and cheating = 0
UNION ALL
select depth, count as count, prog_l, 'zero' as wr_decay
    from wiring where wr_nearby = 0 and cheating = 0
order by wr_decay, depth
;

select * from counthist;




