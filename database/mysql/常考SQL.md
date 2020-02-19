
## 去重distinct
2020.01.15 顺丰同城 golang基础架构面试SQL题

有两个表，Student表中有三个字段：学号id，姓名name；课程Course表有三个字段：学号id，课程编号num，分数。
要求给出学号，姓名，排名的SQL语句，当成绩相同时名次一样，且名次连续，如98分第1名，98并列第1名，97第2名而不是第3名。

分析：
要求名次连续，可以使用distinct先对score去重，得到中间结果：select distinct(score) as ds from Score order by score desc

select Student.id, Student.name, count(ds in
    (select distinct(score) as ds from Score order by score desc) as rank
)
from Student
left join Score on Student.id = Course.id
group by Course.score 
order by Course.score