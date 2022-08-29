#!bash
# shellcheck disable=SC1090
# ~/.bashrc: executed by bash(1) for non-login shells.
# see /usr/share/doc/bash/examples/startup-files (in the package bash-doc)
# for examples

case $- in
*i*) ;; # interactive
*) return ;; 
esac

# If not running interactively, don't do anything
[ -z "$PS1" ] && return


# ------------------------- distro detection -------------------------
export DISTRO
[[ $(uname -r) =~ Microsoft ]] && DISTRO=WSL2 #TODO distinguish WSL1


# ----------------------- environment variables ----------------------
#                           (also see envx)

export USER="${USER:-$(whoami)}"
export GITUSER="$USER"
export REPOS="$HOME/Repos"
export GHREPOS="$REPOS/github.com/$GITUSER"
export DOTFILES="$GHREPOS/dot"
export SCRIPTS="$DOTFILES/scripts"
export SNIPPETS="$DOTFILES/snippets"
export HELP_BROWSER=lynx
export DESKTOP="$HOME/Desktop"
export DOCUMENTS="$HOME/Documents"
export DOWNLOADS="$HOME/Downloads"
export TEMPLATES="$HOME/Templates"
export PUBLIC="$HOME/Public"
export PRIVATE="$HOME/Private"
export PICTURES="$HOME/Pictures"
export MUSIC="$HOME/Music"
export VIDEOS="$HOME/Videos"
export PDFS="$HOME/usb/pdfs"
export VIRTUALMACHINES="$HOME/VirtualMachines"
export WORKSPACES="$HOME/Workspaces" # container home dirs for mounting
export ZETDIR="$GHREPOS/zet"
export ZETTELCASTS="$VIDEOS/ZettelCasts"
export CLIP_DIR="$VIDEOS/Clips"
export CLIP_DATA="$GHREPOS/cmd-clip/data"
export CLIP_VOLUME=0
export CLIP_SCREEN=0
export TERM=xterm-256color
export HRULEWIDTH=73
export EDITOR=vi
export VISUAL=vi
export EDITOR_PREFIX=vi


export LESS_TERMCAP_mb="[35m" # magenta
export LESS_TERMCAP_md="[33m" # yellow
export LESS_TERMCAP_me="" # "0m"
export LESS_TERMCAP_se="" # "0m"
export LESS_TERMCAP_so="[34m" # blue
export LESS_TERMCAP_ue="" # "0m"
export LESS_TERMCAP_us="[4m"  # underline


# ------------------------------ cdpath ------------------------------

export CDPATH=".:$GHREPOS:$DOTFILES:$REPOS:/media/$USER:$HOME"

export GOPATH=$HOME/go

LFCD="$GOPATH/src/github.com/gokcehan/lf/etc/lfcd.sh"  # source
LFCD="/path/to/lfcd.sh"                                #  pre-built binary, make sure to use absolute path
if [ -f "$LFCD" ]; then
    source "$LFCD"
fi

source /root/bash-wakatime/bash-wakatime.sh

export PAGER=bat


# path variable manipulation
export PATH="$PATH:$GOPATH/bin:/etc/lynx:/root/Repos/github.com/root/dot/scripts:/root/.cowsay/cowfiles"

# ------------------------ bash shell options ------------------------

shopt -s checkwinsize
shopt -s expand_aliases
shopt -s globstar
shopt -s dotglob
shopt -s extglob

# ... or force ignoredups and ignorespace
export HISTCONTROL=ignoreboth

# append to the history file, don't overwrite it
shopt -s histappend

# for setting history length see HISTSIZE and HISTFILESIZE in bash(1)
export HISTSIZE=5000
export HISTFILESIZE=10000

set -o vi

# make less more friendly for non-text input files, see lesspipe(1)
[ -x /usr/bin/lesspipe ] && eval "$(SHELL=/bin/sh lesspipe)"

# set variable identifying the chroot you work in (used in the prompt below)
if [ -z "$debian_chroot" ] && [ -r /etc/debian_chroot ]; then
    debian_chroot=$(cat /etc/debian_chroot)
fi

# set a fancy prompt (non-color, unless we know we "want" color)
case "$TERM" in
    xterm-color) color_prompt=yes;;
esac

# uncomment for a colored prompt, if the terminal has the capability; turned
# off by default to not distract the user: the focus in a terminal window
# should be on the output of commands, not on the prompt
#force_color_prompt=yes

if [ -n "$force_color_prompt" ]; then
    if [ -x /usr/bin/tput ] && tput setaf 1 >&/dev/null; then
	# We have color support; assume it's compliant with Ecma-48
	# (ISO/IEC-6429). (Lack of such support is extremely rare, and such
	# a case would tend to support setf rather than setaf.)
	color_prompt=yes
    else
	color_prompt=
    fi
fi

# export PS1='\[\033[33m\][\t]\[\033[00m\] \[\033[32m\]\w$\[\033[00m\] '
# --------------------------- smart prompt ---------------------------
#                 (keeping in bashrc for portability)

PROMPT_LONG=20
PROMPT_MAX=95
PROMPT_AT=@

__ps1() {
	  local P='$' dir="${PWD##*/}" B countme short long double\
	      r='\[\e[31m\]' g='\[\e[30m\]' h='\[\e[34m\]' \
	      u='\[\e[33m\]' p='\[\e[34m\]' w='\[\e[35m\]' \
	      b='\[\e[36m\]' x='\[\e[0m\]'

	    [[ $EUID == 0 ]] && P='#' && u=$r && p=$u # root
	    [[ $PWD = / ]] && dir=/
	    [[ $PWD = "$HOME" ]] && dir='~'

	    B=$(git branch --show-current 2>/dev/null)
	    [[ $dir = "$B" ]] && B=.
	    countme="$USER$PROMPT_AT$(hostname):$dir($B)\$ "

	    [[ $B = master || $B = main ]] && b="$r"
	    [[ -n "$B" ]] && B="$g($b$B$g)"

	    short="$u\u$g$PROMPT_AT$h\h$g:$w$dir$B$p$P$x "
	    long="$g╔ $u\u$g$PROMPT_AT$h\h$g:$w$dir$B\n$g╚ $p$P$x "
	    double="$g╔ $u\u$g$PROMPT_AT$h\h$g:$w$dir\n$g║ $B\n$g╚ $p$P$x "

	    if (( ${#countme} > PROMPT_MAX )); then
	      PS1="$double"
	    elif (( ${#countme} > PROMPT_LONG )); then
	      PS1="$long"
	    else
	      PS1="$short"
	    fi
}

PROMPT_COMMAND="__ps1"

unset color_prompt force_color_prompt

# enable color support of ls and also add handy aliases
if [ -x /usr/bin/dircolors ]; then
    test -r ~/.dircolors && eval "$(dircolors -b ~/.dircolors)" || eval "$(dircolors -b)"
fi

# ------------------------------ aliases -----------------------------
#      (use exec scripts instead, which work from vim and subprocs)

alias cat='bat'
alias find='fdfind'
alias grep='grepall'
alias path='echo "$PATH" \| tr ':' '\n''

alias l='exa'
alias la='exa -a'
alias ls='exa --color=auto'
alias ll='exa -lah'
alias tree='tree -Csuh'

alias h='function hdi(){ howdoi $* -c -n 2; }; hdi'
alias hcook='function hcook(){ HOWDOI_URL=cooking.stackexchange.com howdoi $* ; }; hcook'
alias '?'=duck
alias '??'=google
alias '???'=wiki
alias 'time??'='now'

alias 'cheat'='cht.sh'
alias 'cheatcpp'='cht.sh --shell cpp'


alias dot='cd $DOTFILES'
alias scripts='cd $SCRIPTS'
alias snippets='cd $SNIPPETS'

alias cdhome='cd /mnt/c/Users/atkum/Desktop/'

alias cdcode='cd /mnt/c/Users/atkum/Code'
alias cdkb='cd /mnt/c/Users/atkum/Code/KnowledgeBase'
alias cdzet='cd /mnt/c/Users/atkum/Code/KnowledgeBase/notes/zettelkasten'
alias cdnotes='cd /mnt/c/Users/atkum/Code/KnowledgeBase/notes'

alias cdlib='cd /mnt/c/Users/atkum/Documents/BooksLibrary'

alias lynx='lynx -accept_all_cookies'

alias diff='diff --color'

alias clear='printf "\e[H\e[2J"'
alias cl='printf "\e[H\e[2J"'
alias c='printf "\e[H\e[2J"'

alias subl='/mnt/c/Program\ Files/Sublime\ Text/subl.exe'
alias start='cmd.exe /C start'
alias fs='lf'

# Estimate file space usage to maximum depth
alias du1="du -d 1"

alias cp="spinner cp"
alias mkcd='_mkcd(){ mkdir "$1"; cd "$1";}; _mkcd'

eval "$(dircolors -p | \
    sed 's/ 4[0-9];/ 01;/; s/;4[0-9];/;01;/g; s/;4[0-9] /;01 /' | \
	    dircolors /dev/stdin)"





export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"  # This loads nvm
[ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"  # This loads nvm bash_completion






envx() {
  local envfile="${1:-"$HOME/.env"}"
  [[ ! -e "$envfile" ]] && echo "$envfile not found" && return 1
  while IFS= read -r line; do
     name=${line%%=*}
     value=${line#*=}
     [[ -z "${name}" || $name =~ ^# ]] && continue
     export "$name"="$value"
  done < "$envfile"
} && export -f envx

[[ -e "$HOME/.env" ]] && envx "$HOME/.env"



#-------------------------------------------------------------------
# Cow-spoken fortunes every time you open a terminal
function cowsayfortune {
    NUMOFCOWS=`cowsay -l | tail -n +2 | wc -w`
    WHICHCOW=$((RANDOM%$NUMOFCOWS+1))
    THISCOW=`cowsay -l | tail -n +2 | sed -e 's/\ /\'$'\n/g' | sed $WHICHCOW'q;d'`

    #echo "Selected cow: ${THISCOW}, from ${WHICHCOW}"
    fortune | cowsay -f $THISCOW -W 100 | lolcat
}

alias rand='cowsayfortune'

animals() {
    figlet -f ANSI\ Regular.flf $@ | lolcat
}

figlet -f slant "rwxvishnu" | lolcat
#------------------------------------------------------------------------
