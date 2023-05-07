# ciscoCompare
用于比较思科5GC配置和给定规则表格的差异





### 使用方法

#### 1.

将思科提供的配置表格打开，将前4个页签分别另存为csv文件

文件名分别设置为：

allUserProfile.csv

userProfileBinding.csv

filter.csv

action.csv

#### 2.

将思科提供的配置文件，修改名称为config.txt



#### 3.

运行main.go



#### 4.

运行结果在data文件夹的report.txt文件中

复制文件内容，粘贴到json在线解析网站https://www.json.cn/中查看



