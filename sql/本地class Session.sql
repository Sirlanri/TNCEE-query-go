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

select count(*) from gaokao.18totaldata where 专业="机械设计制造及其自动化" and 成绩=520 
select count(*) from gaokao.18totaldata where 专业="机械设计制造及其自动化"

select 专业, max(成绩) from gaokao.17totaldata where 科类代码='理工' group by 专业 

update gaokao.17lg set maxscore =? where name=?

select count(*) from gaokao.17totaldata where 专业='软件工程' and 成绩 between 510 and 590
select 成绩,count(*) from gaokao.17totaldata where 专业='软件工程' and 成绩 between 505 and 590 group by 成绩 order by 成绩

select name, averank from gaokao.lg19 where averank between 80000 and 90000

select 理工科类累计人数 from gaokao.rank19 where 成绩分数段=520

update gaokao.lg17 set maxrank=? where name=?

select 理工科类累计人数 from gaokao.rank19 where 成绩分数段=527