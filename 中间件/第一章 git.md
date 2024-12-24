## 基本概念

Git 中有四个主要的工作区域（或者说四个状态），它们是工作区（Working Directory）、暂存区（Index）、本地仓库（Local Repository）、远程仓库（Remote Repository）。

1. **工作区（Working Directory）**：
   - 工作区是你当前工作的目录，包含项目的实际文件。在这里，你修改文件、添加新文件或者删除文件。工作区就是你在电脑里能看到的项目目录。
2. **暂存区（Index）**：
   - 暂存区是 Git 的一个概念，也可以称为索引（Index）。它是一个临时存放你的改动的地方，你可以将要提交的修改放在这里，然后一次性提交到本地仓库。通过 `git add` 命令可以把工作区的修改提交到暂存区。
3. **本地仓库（Local Repository）**：
   - 本地仓库是 Git 保存项目历史记录的地方。它包含了完整的 Git 版本库，包括所有的分支、提交历史等。本地仓库通常位于你的项目目录中的 `.git` 目录下。
4. **远程仓库（Remote Repository）**：
   - 远程仓库是存放在网络上的项目仓库，比如托管在 GitHub、GitLab 或 Bitbucket 等平台上。它可以被多个开发者访问和共享。通过 `git push` 和 `git pull` 命令，你可以把本地仓库的修改推送到远程仓库，或者从远程仓库拉取更新到本地仓库。

### 工作流程示例：

- 当你在工作区修改了文件后，使用 `git add` 将修改添加到暂存区。
- 使用 `git commit` 将暂存区的修改提交到本地仓库，生成一个新的提交。
- 使用 `git push` 将本地仓库的修改推送到远程仓库。
- 使用 `git pull` 从远程仓库拉取更新到本地仓库。



## git使用

### 步骤 1：安装 Git

首先，你需要在你的计算机上安装 Git。Git 支持多个操作系统，包括 Windows、Mac 和 Linux。你可以从 [Git 官网](https://git-scm.com/) 下载安装程序，并按照提示进行安装。

### 步骤 2：配置 Git

安装完成后，打开命令行界面（如终端或命令提示符），配置 Git 的全局用户信息，包括用户名和邮箱地址。这些信息将会出现在你的提交记录中，方便识别是谁做了哪些修改。

```
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"
```

### 步骤 3：创建一个新的 Git 仓库

现在，让我们来创建一个新的 Git 仓库。假设你有一个项目目录，进入这个目录并执行以下命令：

```
cd your-project-directory
git init
```

这会在你的项目目录中初始化一个新的 Git 仓库，Git 会在该目录下生成一个 `.git` 的子目录，用于存储版本控制相关的信息。

### 步骤 4：进行基本的 Git 操作

现在，你已经有了一个空的 Git 仓库，可以开始进行基本的操作了：

- **查看仓库状态**：可以随时查看工作区和暂存区的状态，看看哪些文件被修改了但还未提交。

  ```
  git status
  ```

- **添加文件到暂存区**：将工作区的修改添加到暂存区，准备提交。

  ```
  git add filename    # 添加单个文件
  git add .           # 添加所有修改的文件
  ```

- **提交到本地仓库**：将暂存区的修改提交到本地仓库。

  ```
  git commit -m "Commit message describing your changes"
  ```

### 步骤 5：查看历史记录和撤销修改

- **查看提交历史**：可以查看所有的提交记录。

  ```
  git log
  ```

- **撤销修改**：可以撤销工作区的修改或者撤销提交。

  ```
  git checkout -- filename    # 撤销工作区的修改
  git reset HEAD filename     # 撤销暂存区的修改
  ```

### 步骤 6：远程仓库的使用

如果你有一个远程仓库（比如 GitHub），可以将本地仓库与之关联，进行远程操作。

- **关联远程仓库**：将本地仓库与远程仓库关联。

  ```
  git remote add origin <remote-repository-url>
  ```

- **推送到远程仓库**：将本地仓库的修改推送到远程仓库。

  ```
  git push -u origin master    # 推送到主分支（master）上
  ```

- **从远程仓库拉取更新**：获取远程仓库的更新到本地仓库。

  ```
  git pull origin master    # 拉取远程主分支（master）的更新
  ```

### 步骤 7：分支操作

Git 的分支功能非常强大，可以帮助你管理不同的功能或修复的工作流。

- **创建分支**：创建一个新的分支来开发新功能或者修复问题。

  ```
  git branch new-feature    # 创建名为 new-feature 的新分支
  ```

- **切换分支**：切换到其他分支进行工作。

  ```
  git checkout new-feature    # 切换到 new-feature 分支
  ```

- **合并分支**：将一个分支的修改合并到当前分支。

  ```
  git merge new-feature    # 将 new-feature 分支的修改合并到当前分支
  ```

## 提交一个新的项目

echo "# yourProgram" >> README.md

git init

git add README.md

git commit -m "first commit"

git branch -M main

git remote add origin https://github.com/Yuyq1114/yourProgram.git

git push -u origin main



## 提交一个以存在的项目

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



##   命令详解

### **1. 初始化仓库**

- **`git init`**
  初始化一个 Git 仓库，将目录变为受 Git 管理的项目。

  ```
  git init
  ```

- **`git clone`**
  从远程仓库克隆项目到本地。

  ```
  git clone <repository_url>
  ```

------

### **2. 工作区相关命令**

#### 查看状态

- `git status`

  查看工作区和暂存区文件状态。

  ```
  git status
  ```

#### 比较文件

- `git diff`

  查看工作区中修改的内容（未暂存的部分）。

  ```
  git diff
  ```

------

### **3. 暂存区相关命令**

#### 添加文件到暂存区

- **`git add <file>`**
  将指定文件添加到暂存区。

  ```
  git add file.txt
  ```

- **`git add .`**
  添加当前目录下的所有更改到暂存区。

#### 从暂存区中移除文件

- `git reset <file>`

  将暂存区中的指定文件移除，但保留工作区的修改。

  ```
  git reset file.txt
  git reset --hard HEAD//如果不想保留任何更改，可以强制重置分支
  ```

------

### **4. 本地仓库相关命令**

#### 提交更改

- **`git commit -m "<message>"`**
  将暂存区的内容提交到本地仓库，并添加提交信息。

  ```
  git commit -m "Initial commit"
  ```

- **`git commit -a -m "<message>"`**
  跳过暂存区，直接提交工作区中已跟踪的文件修改。

#### 修改提交

- **`git commit --amend`**
  修改上一次提交的信息或文件。

#### fetch

- 获取某个远程仓库（例如 `origin`）的更新：

  ```
  git fetch origin
  ```

#### 查看提交历史

- **`git log`**
  查看提交历史。

  ```
  git log
  ```

- **`git log --oneline`**
  简化输出的提交历史。

  ```
  git log --oneline
  ```

------

### **5. 远程仓库相关命令**

#### 添加远程仓库

- `git remote add <name> <url>`

  将远程仓库关联到本地仓库。

  ```
  git remote add origin https://github.com/user/repo.git
  ```

#### 推送更改

- `git push <remote> <branch>`

  将本地分支内容推送到远程分支。

  ```
  git push origin main
  ```

#### 拉取更改

- `git pull <remote> <branch>`

  拉取远程分支的最新内容并合并到本地。

  ```
  git pull origin main
  ```

#### 查看远程仓库

- `git remote -v`

  查看远程仓库的详细信息。

  ```
  git remote -v
  ```

------

### **6. 分支管理**

#### 查看分支

- **`git branch`**
  查看本地分支列表。

  ```
  git branch
  ```

- **`git branch -r`**
  查看远程分支列表。

#### 创建分支

- `git branch <branch_name>`

  创建新分支。

  ```
  git branch feature-1
  ```

#### 切换分支

- **`git checkout <branch_name>`**
  切换到指定分支。

  ```
  git checkout main
  ```

- **`git switch <branch_name>`**
  使用较新的切换分支方式。

#### 删除分支

- `git branch -d <branch_name>`

  删除本地分支。

  ```
  git branch -d feature-1
  ```

------

#### 合并分支

##### **1. 合并分支的基本概念**

- **目标分支（当前分支）**：需要合并更改到的分支。
- **源分支**：提供更改的分支。
- 合并后，目标分支将包含源分支的更改记录。

------

#### **2. 使用命令行合并分支**

##### **（1）检查当前分支**

确保切换到需要合并更改的目标分支：

```
git checkout <target-branch>
```

示例：

```
git checkout main
```

##### **（2）执行合并**

运行以下命令，将源分支合并到当前分支：

```
git merge <source-branch>
```

示例：

```
git merge feature-branch
```

##### **（3）处理合并冲突（如果有）**

- 如果 Git 检测到冲突，合并将暂停，你会看到类似下面的信息：

  ```
  CONFLICT (content): Merge conflict in <file>
  ```
  
- 解决冲突后，运行以下命令完成合并：

  ```
  git add <conflicted-file>
  git commit
  ```

### **7. 标签管理**

#### 创建标签

- **`git tag <tag_name>`**
  创建轻量标签。

  ```
  git tag v1.0
  ```

- **`git tag -a <tag_name> -m "<message>"`**
  创建带注释的标签。

  ```
  git tag -a v1.0 -m "Release version 1.0"
  ```

#### 推送标签

- `git push <remote> <tag_name>`

  将标签推送到远程仓库。

  ```
  git push origin v1.0
  ```

------

### **8. 撤销操作**

#### 恢复文件

- `git checkout -- <file>`

  恢复工作区中文件的修改。

  ```
  git checkout -- file.txt
  ```

#### 回滚提交

- `git reset`

  回滚到指定的提交：

  - **`--soft`**: 保留提交和暂存区内容。
  - **`--mixed`**: 保留提交，但清空暂存区。
  - **`--hard`**: 清空提交和暂存区，直接恢复到指定状态。
  
- `git revert`：通过生成一个新的提交来撤销某个提交，不修改提交历史，适合公开历史的项目。

------

### **9. 其他命令**

#### 清理未跟踪文件

- `git clean -f`

  删除未跟踪的文件。

  ```
  git clean -f
  ```

#### 查看变更

- `git show`

  查看最新一次提交的详细变更。

  ```
  git show
  ```

#### 比较版本

- `git diff <commit1> <commit2>`

  比较两个提交的差异。

  ```
  git diff HEAD~1 HEAD
  ```