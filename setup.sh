#!/usr/bin/env bash

# This install script is useful if running on a RHEL install or image. If you're
# using another OS or Linux distro, this script can be used as a guide to ensure
# that all dependencies are satisfied.

go_version="1.9.2";
ruby_version="2.4.0";
if [ $# -ne 2 ]; then
    echo "No/invalid parameters specified. Exiting";
    exit 1;
else
    go_version=$1;
    ruby_version=$2;
    echo "Using supplied version params:";
    echo "    Go: $go_version";
    echo "    Ruby: $ruby_version";
fi

home="/home/ec2-user/";

# switch to superuser
sudo su;

# update local packages
yum update -y -q -e 0;
yum install wget -y -q -e 0;
echo "Finished updating system packages.";

# setup system
cd $home;
mkdir -p ./go;

# install golang
go_file="go1.9.2.linux-amd64.tar.gz";
wget "https://storage.googleapis.com/golang/$go_file" -q;
tar -C /usr/local -xzf $go_file;
go_result=$( go version );
if [[ $go_result == "go version"* ]]; then
    echo -e "\e[32mGolang version $go_version successfully installed.\e[0m";
    go get -u github.com/ggrpmnn/orion;
    go get -u github.com/GoAST/gas;
else
    echo -e "\e[31mGolang install failed; please try again manually.\e[0m";
fi
rm $go_file;

# install ruby (and gems)
ruby_file="ruby-2.4.0.tar.gz";
wget "https://cache.ruby-lang.org/pub/ruby/$ruby_file" -q;
tar -C /usr/local -xzf $ruby_file;
cd /usr/local/$ruby_version;
./configure;
make --silent;
make install --silent;
ruby_result=$( ./ruby -v );
if [[ $ruby_result == "ruby"* ]]; then
    echo -e "\e[32mRuby version $ruby_version successfully installed.\e[0m";
    gem update --system;
    gem install bundler;
else
    echo -e "\e[31mRuby install failed; please try again manually.\e[0m";
fi
cd $home;
rm $ruby_file;

# add preferred vals to setup profile (optional)
echo "" >> $home/.bash_profile;
echo "# for Golang" >> $home/.bash_profile;
echo "export GOPATH=/home/ec2-user/go" >> $home/.bash_profile;
echo "export GOBIN=/home/ec2-user/go/bin" >> $home/.bash_profile;
echo "export PATH=$PATH:$GOBIN:/usr/local/go/bin" >> $home/.bash_profile;
echo "" >> $home/.bash_profile;
echo "# for Ruby" >> $home/.bash_profile;
echo "export PATH=$PATH:/usr/local/ruby-$ruby_version:/usr/local/ruby-$ruby_version/lib:/usr/local/ruby-$ruby_version/bin">> $home/.bash_profile;
echo "" >> $home/.bash_profile;
echo "# on login" >> $home/.bash_profile;
# this file should be added later by hand to keep creds safe
echo "source /home/ec2-user/go/src/github.com/ggrpmnn/orion/.ign/env" >> $home/.bash_profile;
echo "cd /home/ec2-user/go/src/github.com/ggrpmnn/orion" >> $home/.bash_profile;
