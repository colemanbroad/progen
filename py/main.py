import pandas
import sqlite3
import plotly.express as px

def plot():
    conn = sqlite3.connect("../wiring.sqlite")
       
    # prog_l wr_decay wr_nearby n_prog depth count cheating

    q = """
        select depth, count/(prog_l) as count, prog_l, wr_decay from wiring where wr_nearby = 1 and cheating = 1
        UNION ALL
        select depth, count/(prog_l) as count, prog_l, 'zero' as wr_decay from wiring where wr_nearby = 0 and cheating = 1
        order by wr_decay, depth
    """
    df = pandas.read_sql_query(q,conn)
    fig = px.line(df, x="depth", y="count", color="prog_l", facet_row="wr_decay")
    fig.show()
    # fig.update_traces(marker=dict(size=3))

plot()
