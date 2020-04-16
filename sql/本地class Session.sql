select 理工科类累计人数 from `gaokao`.`rank18` 
    where 成绩分数段 in(
        select 平均分 from `gaokao`.`lg` 
            where 年份=2018
    )

select 专业,count(性别) from gaokao.17totaldata where 性别="男" group by 专业;
select 专业,count(性别) from gaokao.17totaldata where 性别='女' and 科类代码='理工' and 专业='' group by 专业

select count(性别) from gaokao.17totaldata where 性别='女' and 科类代码='理工' and 专业="机械设计制造及其自动化" group by 专业

select minscore,minrank,avescore,averank from gaokao.lg17 where name='工业设计'
union
select minscore,minrank,avescore,averank from gaokao.lg18 where name='工业设计'
union
select minscore,minrank,avescore,averank from gaokao.lg19 where name='工业设计'