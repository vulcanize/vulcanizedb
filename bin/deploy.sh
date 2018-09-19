if [ $TRAVIS_BRANCH == 'staging' ]; then
  sup --debug staging deploy
fi
