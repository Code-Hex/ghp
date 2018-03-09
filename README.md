# ghp(1)

Create a new project on the [ghq](https://github.com/motemen/ghq) root

# Synopsis
    
    ghp newProjectName       # It will make directory on ghq root and git init
    cd $(ghp newProjectName) # easy to change directory

or

    ghp --with-license newProjectName

![c67ac587f651b43dbd3c57fbb0dc4833](https://user-images.githubusercontent.com/6500104/37215511-76ac7cac-23fb-11e8-84ea-494d3d6e4fd0.gif)

also you can use combination with [peco](https://github.com/peco/peco) and ghq. see [sample usage](https://github.com/peco/peco/wiki/Sample-Usage#pecoghq--ghq--peco-miyagawa)

# Installation

    go get github.com/Code-Hex/ghp/cmd/ghp
