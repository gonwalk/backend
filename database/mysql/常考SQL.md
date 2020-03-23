
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

# LeetCode 176. 第二高的薪水

编写一个 SQL 查询，获取 Employee 表中第二高的薪水（Salary） 。

+----+--------+
| Id | Salary |
+----+--------+
| 1  | 100    |
| 2  | 200    |
| 3  | 300    |
+----+--------+
例如上述 Employee 表，SQL查询应该返回 200 作为第二高的薪水。如果不存在第二高的薪水，那么查询应返回 null。

+---------------------+
| SecondHighestSalary |
+---------------------+
| 200                 |
+---------------------+

```sql
# 方式一：使用中间表
select (
    select distinct Salary 
    from Employee 
    order by Salary desc
    limit 1 offset 1
) as SecondHighestSalary

# 方式二：使用IFNULL
-- SELECT
--     IFNULL(
--       (SELECT DISTINCT Salary
--        FROM Employee
--        ORDER BY Salary DESC
--         LIMIT 1 OFFSET 1),
--     NULL) AS SecondHighestSalary
```


# LeetCode 183. 从不订购的客户


某网站包含两个表，Customers 表和 Orders 表。编写一个 SQL 查询，找出所有从不订购任何东西的客户。

Customers 表：

+----+-------+
| Id | Name  |
+----+-------+
| 1  | Joe   |
| 2  | Henry |
| 3  | Sam   |
| 4  | Max   |
+----+-------+
Orders 表：

+----+------------+
| Id | CustomerId |
+----+------------+
| 1  | 3          |
| 2  | 1          |
+----+------------+
例如给定上述表格，你的查询应返回：

+-----------+
| Customers |
+-----------+
| Henry     |
| Max       |
+-----------+

```sql
# Write your MySQL query statement below
select Customers.Name as Customers      # 注意：这里不能使用distinct，否则会出错，因为Customers表中可能有重名的，但是Id不同的记录
from Customers
where Customers.Id not in (
    select Orders.CustomerId
    from Orders
)
```

# LeetCode 196. 删除重复的电子邮箱

题目来源：https://leetcode-cn.com/problems/delete-duplicate-emails/

## 题目描述

编写一个 SQL 查询，来删除 Person 表中所有重复的电子邮箱，重复的邮箱里只保留 Id 最小 的那个。

+----+------------------+
| Id | Email            |
+----+------------------+
| 1  | john@example.com |
| 2  | bob@example.com  |
| 3  | john@example.com |
+----+------------------+
Id 是这个表的主键。
例如，在运行你的查询语句之后，上面的 Person 表应返回以下几行:

+----+------------------+
| Id | Email            |
+----+------------------+
| 1  | john@example.com |
| 2  | bob@example.com  |
+----+------------------+
 

提示：

执行 SQL 之后，输出是整个 Person 表。
使用 delete 语句。


## SQL语句

我们可以使用以下代码，将此表与它自身在电子邮箱列中连接起来。

SELECT p1.*
FROM Person p1,
    Person p2
WHERE
    p1.Email = p2.Email
;
然后我们需要找到其他记录中具有相同电子邮件地址的更大 ID。所以我们可以像这样给 WHERE 子句添加一个新的条件。

SELECT p1.*
FROM Person p1,
    Person p2
WHERE
    p1.Email = p2.Email AND p1.Id > p2.Id
;
因为我们已经得到了要删除的记录，所以我们最终可以将该语句更改为 DELETE。
```sql
DELETE p1 FROM Person p1,
    Person p2
WHERE
    p1.Email = p2.Email AND p1.Id > p2.Id
```

# LeetCode 197. 上升的温度

题目来源：https://leetcode-cn.com/problems/rising-temperature/

知识点：记住两个时间函数——datediff和timestampdiff
```
datediff(日期1, 日期2)：得到的结果是日期1与日期2相差的天数。如果日期1比日期2大，结果为正；如果日期1比日期2小，结果为负。

timestampdiff(时间类型, 日期1, 日期2)：这个函数和上面datediff的正、负号规则刚好相反。日期1大于日期2，结果为负，日期1小于日期2，结果为正。
在“时间类型”的参数位置，通过添加DAY, HOUR, 等关键词，来规定计算天数差、小时数差、还是分钟数差。
```

## 题目描述

给定一个 Weather 表，编写一个 SQL 查询，来查找与之前（昨天的）日期相比温度更高的所有日期的 Id。

+---------+------------------+------------------+
| Id(INT) | RecordDate(DATE) | Temperature(INT) |
+---------+------------------+------------------+
|       1 |       2015-01-01 |               10 |
|       2 |       2015-01-02 |               25 |
|       3 |       2015-01-03 |               20 |
|       4 |       2015-01-04 |               30 |
+---------+------------------+------------------+
例如，根据上述给定的 Weather 表格，返回如下 Id:

+----+
| Id |
+----+
|  2 |
|  4 |
+----+

## SQL

TIMESTAMPDIFF能干什么，可以计算相差天数、小时、分钟和秒，相比于datediff函数要灵活很多。格式是时间小的前，时间大的放在后面。
计算相差天数：

select TIMESTAMPDIFF(DAY,'2019-05-20', '2019-05-21'); # 1
计算相差小时数：

select TIMESTAMPDIFF(HOUR, '2015-03-22 07:00:00', '2015-03-22 18:00:00'); # 11
计算相差秒数：

select TIMESTAMPDIFF(SECOND, '2015-03-22 07:00:00', '2015-03-22 7:01:01'); # 61
采用联结表的方式，条件是：1）与之前的日期相差为 1天，2）当天温度比之前一天的温度高
```sql
select w1.Id
from Weather as w1, Weather as w2
where TIMESTAMPDIFF(DAY, w2.RecordDate, w1.RecordDate) = 1 AND w1.Temperature > w2.Temperature
```