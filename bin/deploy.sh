if [ $TRAVIS_BRANCH == 'staging' ]; then
  sup --debug staging deploy
elif [ $TRAVIS_BRANCH == 'master' ]; then
  sup --debug prod deploy
fi
