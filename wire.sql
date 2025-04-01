.mode col

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
