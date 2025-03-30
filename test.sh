
read -r -d '' query << EOM
select depth, count/(prog_l) as count, prog_l, wr_decay from wiring where wr_nearby = 1 and cheating = 1
UNION ALL
select depth, count/(prog_l) as count, prog_l, 'zero' as wr_decay from wiring where wr_nearby = 0 and cheating = 1
order by wr_decay, depth
EOM

# sqlite3 wiring.sqlite << EOM
# .mode col
# select depth, count/(prog_l) as count, prog_l, wr_decay from wiring where wr_nearby = 1 and cheating = 1
# UNION ALL
# select depth, count/(prog_l) as count, prog_l, 'zero' as wr_decay from wiring where wr_nearby = 0 and cheating = 1
# order by wr_decay, depth
# limit 20
# EOM

sqlpeek -d wiring.sqlite -q $query -x xyl

# df = pandas.read_sql_query(q,conn)
# fig = px.line(df, x="depth", y="count", color="prog_l", facet_row="wr_decay")
