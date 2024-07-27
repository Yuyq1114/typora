# 提交一个新的项目

echo "# yourProgram" >> README.md

git init

git add README.md

git commit -m "first commit"

git branch -M main

git remote add origin https://github.com/Yuyq1114/yourProgram.git

git push -u origin main



# 提交一个以存在的项目

<!--git remote add origin https://github.com/Yuyq1114/yourProgram.git-->

<!--git branch -M main-->

<!--git push -u origin main-->

git  add fileName

git commit -m main

git push origin main

**遇到冲突**

git pull origin branch_name

git add conflicted_file

git commit -m "Resolved merge conflict in conflicted_file"

git push origin branch_name